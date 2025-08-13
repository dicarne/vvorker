package alioss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"vvorker/common"
	"vvorker/conf"

	aoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/gin-gonic/gin"
)

// 从请求头中获取参数并创建客户端
func getOSSClient(c *gin.Context) (*aoss.Client, error) {
	resourceID := c.GetHeader("ResourceID")
	endpoint := c.GetHeader("Endpoint")
	accessKeyID := c.GetHeader("AccessKeyID")
	secretAccessKey := c.GetHeader("SecretAccessKey")
	if conf.AppConfigInstance.MinioSingleBucketMode {
		accessKeyID = conf.AppConfigInstance.ServerMinioAccess
		secretAccessKey = conf.AppConfigInstance.ServerMinioSecret
	}
	region := c.GetHeader("Region")
	useSSLStr := c.GetHeader("UseSSL")
	if len(resourceID) != 0 {
		endpoint = fmt.Sprintf("%s:%d", conf.AppConfigInstance.ServerMinioHost, conf.AppConfigInstance.ServerMinioPort)
		useSSLStr = fmt.Sprintf("%v", conf.AppConfigInstance.ServerMinioUseSSL)
	}

	// logrus.Infof("Endpoint: %s, AccessKeyID: %s, SecretAccessKey: %s, UseSSL: %s, Region: %s",
	// 	endpoint, accessKeyID, secretAccessKey, useSSLStr, region)

	useSSL, err := strconv.ParseBool(useSSLStr)
	if err != nil {
		return nil, err
	}

	cfg := aoss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithDisableSSL(!useSSL).
		WithUsePathStyle(conf.AppConfigInstance.ServerMinioBucketLoopUp == 2)

	client := aoss.NewClient(cfg)
	return client, nil
}

func getMinioConfig(c *gin.Context) (string, string) {
	bucketName := c.GetHeader("Bucket")
	objectName := c.GetHeader("Object")
	if conf.AppConfigInstance.MinioSingleBucketMode {
		if objectName[0] == '/' {
			objectName = objectName[1:]
		}
		bucketName = conf.AppConfigInstance.MinioSingleBucketName
		objectName = c.GetHeader("ResourceID") + "/" + objectName
		if conf.AppConfigInstance.MinioSingleBucketPrefix != "" {
			objectName = conf.AppConfigInstance.MinioSingleBucketPrefix + "/" + objectName
		}
	}

	return bucketName, objectName
}

// DownloadFile 下载文件接口
func DownloadFile(c *gin.Context) {
	mc, err := getOSSClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}

	bucketName, objectName := getMinioConfig(c)

	// 使用 GetObject 获取文件流
	obj, err := mc.GetObject(context.Background(), &aoss.GetObjectRequest{
		Bucket: &bucketName,
		Key:    &objectName,
	})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get object", gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", *obj.ContentType)
	c.Header("Content-Length", strconv.FormatInt(obj.ContentLength, 10))
	c.Header("Content-Disposition", "attachment; filename="+objectName)

	// 将文件流写入响应体
	_, err = io.Copy(c.Writer, obj.Body)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to copy object to response", gin.H{"error": err.Error()})
		return
	}

	// 确保所有数据都被刷新到客户端
	c.Writer.Flush()
}

// UploadFile 上传文件接口
func UploadFile(c *gin.Context) {
	mc, err := getOSSClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}

	bucketName, objectName := getMinioConfig(c)

	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to get file from form", gin.H{"error": err.Error()})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to open file", gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// 上传文件到 MinIO
	info, err := mc.PutObject(context.Background(), &aoss.PutObjectRequest{
		Bucket: &bucketName,
		Key:    &objectName,
		Body:   src,
	})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to upload file", gin.H{"error": err.Error()})
		return
	}

	common.RespOK(c, "File uploaded successfully", gin.H{"info": info})
}

package alioss1

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"vvorker/common"
	"vvorker/conf"

	aoss "github.com/aliyun/aliyun-oss-go-sdk/oss"
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

	clientOptions := []aoss.ClientOption{}
	clientOptions = append(clientOptions, aoss.Region(region))
	switch conf.AppConfigInstance.ServerOSSAuthVersion {
	case 4:
		clientOptions = append(clientOptions, aoss.AuthVersion(aoss.AuthV4))
	case 2:
		clientOptions = append(clientOptions, aoss.AuthVersion(aoss.AuthV2))
	case 1:
		clientOptions = append(clientOptions, aoss.AuthVersion(aoss.AuthV1))
	default:
		clientOptions = append(clientOptions, aoss.AuthVersion(aoss.AuthV4))
	}

	if useSSL {
		endpoint = "https://" + endpoint
	} else {
		endpoint = "http://" + endpoint
	}

	client, err := aoss.New(endpoint, accessKeyID, secretAccessKey, clientOptions...)
	if err != nil {
		return nil, err
	}
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
	client, err := getOSSClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}

	bucketName, objectName := getMinioConfig(c)

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get bucket", gin.H{"error": err.Error()})
		return
	}

	// 下载文件到流。
	body, err := bucket.GetObject(objectName)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get object", gin.H{"error": err.Error()})
		return
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	// 将文件流写入响应体
	_, err = io.Copy(c.Writer, body)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to copy object to response", gin.H{"error": err.Error()})
		return
	}

	// 确保所有数据都被刷新到客户端
	c.Writer.Flush()
}

// UploadFile 上传文件接口
func UploadFile(c *gin.Context) {
	client, err := getOSSClient(c)
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

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get bucket", gin.H{"error": err.Error()})
		return
	}

	err = bucket.PutObject(objectName, src)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to put object", gin.H{"error": err.Error()})
		return
	}

	common.RespOK(c, "File uploaded successfully", gin.H{"info": "success"})
}

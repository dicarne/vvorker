package oss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"vorker/conf"
	"vorker/entities"
	"vorker/models"
	"vorker/utils"
	"vorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// MinioClient 定义Minio客户端结构体
type MinioClient struct {
	Client *minio.Client
}

// NewMinioClient 创建新的Minio客户端
func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool, region string) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	return &MinioClient{Client: client}, nil
}

// getMinioClient 从请求头中获取参数并创建 Minio 客户端
func getMinioClient(c *gin.Context) (*MinioClient, error) {
	resourceID := c.GetHeader("ResourceID")
	endpoint := c.GetHeader("Endpoint")
	accessKeyID := c.GetHeader("AccessKeyID")
	secretAccessKey := c.GetHeader("SecretAccessKey")
	region := c.GetHeader("Region")
	useSSLStr := c.GetHeader("UseSSL")
	if len(resourceID) != 0 {
		endpoint = fmt.Sprintf("%s:%d", conf.AppConfigInstance.ServerMinioHost, conf.AppConfigInstance.ServerMinioPort)
		useSSLStr = fmt.Sprintf("%v", conf.AppConfigInstance.ServerMinioUseSSL)
	}

	logrus.Infof("Endpoint: %s, AccessKeyID: %s, SecretAccessKey: %s, UseSSL: %s, Region: %s",
		endpoint, accessKeyID, secretAccessKey, useSSLStr, region)

	useSSL, err := strconv.ParseBool(useSSLStr)
	if err != nil {
		return nil, err
	}

	return NewMinioClient(endpoint, accessKeyID, secretAccessKey, useSSL, region)
}

// DownloadFile 下载文件接口
func DownloadFile(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
		return
	}

	bucketName := c.GetHeader("Bucket")
	objectName := c.GetHeader("Object")

	// 使用 GetObject 获取文件流
	obj, err := mc.Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer obj.Close()

	// 获取文件信息，用于设置响应头
	stat, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))
	c.Header("Content-Disposition", "attachment; filename="+objectName)

	// 将文件流写入响应体
	_, err = io.Copy(c.Writer, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保所有数据都被刷新到客户端
	c.Writer.Flush()
}

// UploadFile 上传文件接口
func UploadFile(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
		return
	}

	bucketName := c.GetHeader("Bucket")
	objectName := c.GetHeader("Object")

	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form: " + err.Error()})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	defer src.Close()

	// 上传文件到 MinIO
	info, err := mc.Client.PutObject(context.Background(), bucketName, objectName, src, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "info": info})
}

// DeleteFile 删除文件接口
func DeleteFile(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
		return
	}

	bucketName := c.GetHeader("Bucket")
	objectName := c.GetHeader("Object")

	err = mc.Client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// ListBuckets 查询bucket接口
func ListBuckets(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
		return
	}

	buckets, err := mc.Client.ListBuckets(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}

func ListObjects(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
		return
	}
	bucketName := c.GetHeader("Bucket")
	objectName := c.GetHeader("Path")
	recursive := c.GetHeader("Recursive")
	// 调用 ListObjects 方法
	objectCh := mc.Client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    objectName,
		Recursive: recursive == "true",
	})
	// 遍历对象并构建文件列表
	var files []string
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
			return
		}
		files = append(files, object.Key)
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

// CreateNewOSSResources 创建新的OSS资源
func CreateNewOSSResourcesEndpoint(c *gin.Context) {
	var req entities.CreateNewResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// valid
	if req.UserID == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID and Name are required"})
		return
	}
	mc, err := NewMinioClient(conf.AppConfigInstance.ServerMinioHost,
		conf.AppConfigInstance.ServerMinioAccess,
		conf.AppConfigInstance.ServerMinioSecret,
		conf.AppConfigInstance.ServerMinioUseSSL,
		conf.AppConfigInstance.ServerMinioRegion)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
		return
	}
	// 将 req.UID 转换为 uint64 类型
	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert UserID to uint64: " + err.Error()})
		return
	}
	resource := models.OSS{
		Name:      req.Name,
		UserID:    userID,
		UID:       utils.GenerateUID(),
		Bucket:    "",
		Region:    conf.AppConfigInstance.ServerMinioRegion,
		AccessKey: "",
		SecretKey: "",
	}
	resource.Bucket = resource.UID

	// 检查 bucket 是否存在，如果不存在则创建
	exists, err := mc.Client.BucketExists(context.Background(), resource.Bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check bucket existence: " + err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusOK, gin.H{"message": "Bucket already exists"})
		return
	}
	err = mc.Client.MakeBucket(context.Background(), resource.Bucket, minio.MakeBucketOptions{Region: conf.AppConfigInstance.ServerMinioRegion})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bucket: " + err.Error()})
		return
	}

	account, err := CreateNewServiceAccount(resource.Bucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service account: " + err.Error()})
		return
	}

	resource.AccessKey = account.AccessKey
	resource.SecretKey = account.SecretKey
	resource.Expiration = account.Expiration

	// 创建新的 OSS 资源
	db := database.GetDB()
	if err := db.Create(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OSS resource: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OSS resource created successfully",
		"name":   resource.Name,
		"uid":    resource.UID,
		"status": 0,
	})

}

// 删除指定OSS资源
func DeleteOSSResourcesEndpoint(c *gin.Context) {
	var req entities.DeleteResourcesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// valid
	if req.UID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID is required"})
		return
	}

	db := database.GetDB()
	var resource models.OSS
	if err := db.Where("uid = ?", req.UID).First(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find OSS resource: " + err.Error()})
		return
	}
	// 删除OSS资源
	if err := db.Delete(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete OSS resource: " + err.Error()})
		return
	}
	// 删除Minio中的Bucket
	// mc, err := NewMinioClient(conf.AppConfigInstance.ServerMinioEndpoint,
	// 	conf.AppConfigInstance.ServerMinioAccess,
	// 	conf.AppConfigInstance.ServerMinioSecret,
	// 	conf.AppConfigInstance.ServerMinioUseSSL,
	// 	conf.AppConfigInstance.ServerMinioRegion)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create Minio client: " + err.Error()})
	// 	return
	// }
	// err = mc.Client.RemoveBucket(context.Background(), resource.Bucket)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bucket: " + err.Error()})
	// 	return
	// }
	DeleteServiceAccount(resource.AccessKey)
	c.JSON(http.StatusOK, gin.H{
		"message": "OSS resource deleted successfully",
		"status":  0,
	})
}

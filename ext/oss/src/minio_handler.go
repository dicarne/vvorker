package oss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

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

	logrus.Infof("Endpoint: %s, AccessKeyID: %s, SecretAccessKey: %s, UseSSL: %s, Region: %s",
		endpoint, accessKeyID, secretAccessKey, useSSLStr, region)

	useSSL, err := strconv.ParseBool(useSSLStr)
	if err != nil {
		return nil, err
	}

	return NewMinioClient(endpoint, accessKeyID, secretAccessKey, useSSL, region)
}

func getMinioConfig(c *gin.Context) (string, string) {
	bucketName := c.GetHeader("Bucket")
	objectName := c.GetHeader("Object")
	if conf.AppConfigInstance.MinioSingleBucketMode {
		bucketName = conf.AppConfigInstance.MinioSingleBucketName
		objectName = bucketName + "/" + objectName
	}

	return bucketName, objectName
}

// DownloadFile 下载文件接口
func DownloadFile(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}

	bucketName, objectName := getMinioConfig(c)

	// 使用 GetObject 获取文件流
	obj, err := mc.Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get object", gin.H{"error": err.Error()})
		return
	}
	defer obj.Close()

	// 获取文件信息，用于设置响应头
	stat, err := obj.Stat()
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to get object stat", gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", strconv.FormatInt(stat.Size, 10))
	c.Header("Content-Disposition", "attachment; filename="+objectName)

	// 将文件流写入响应体
	_, err = io.Copy(c.Writer, obj)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to copy object to response", gin.H{"error": err.Error()})
		return
	}

	// 确保所有数据都被刷新到客户端
	c.Writer.Flush()
}

// InitiateMultipartUpload 初始化分块上传
func InitiateMultipartUpload(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}
	coreClient := &minio.Core{Client: mc.Client}

	bucketName, objectName := getMinioConfig(c)

	uploadID, err := coreClient.NewMultipartUpload(context.Background(), bucketName, objectName, minio.PutObjectOptions{})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to initiate multipart upload", gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"UploadId": uploadID,
	})
}

// UploadPart 上传分块
func UploadPart(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}
	coreClient := &minio.Core{Client: mc.Client}

	bucketName, objectName := getMinioConfig(c)
	uploadID := c.GetHeader("x-amz-upload-id")
	partNumberStr := c.GetHeader("x-amz-part-number")
	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Invalid part number", gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to get file from form", gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	part, err := coreClient.PutObjectPart(context.Background(), bucketName, objectName, uploadID, partNumber, file, header.Size, minio.PutObjectPartOptions{})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to upload part", gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ETag": part.ETag,
	})
}

// CompleteMultipartUpload 完成分块上传
func CompleteMultipartUpload(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}
	coreClient := &minio.Core{Client: mc.Client}

	bucketName, objectName := getMinioConfig(c)
	uploadID := c.GetHeader("x-amz-upload-id")

	var completeRequest struct {
		Parts []minio.CompletePart `json:"Parts"`
	}

	if err := c.ShouldBindJSON(&completeRequest); err != nil {
		common.RespErr(c, http.StatusBadRequest, "Invalid parts data", gin.H{"error": err.Error()})
		return
	}

	uploadInfo, err := coreClient.CompleteMultipartUpload(context.Background(), bucketName, objectName, uploadID, completeRequest.Parts, minio.PutObjectOptions{})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to complete multipart upload", gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload completed successfully",
		"Location": uploadInfo.Location,
		"Bucket":   uploadInfo.Bucket,
		"Key":      uploadInfo.Key,
		"ETag":     uploadInfo.ETag,
	})
}

// AbortMultipartUpload 中断分块上传
func AbortMultipartUpload(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}
	coreClient := &minio.Core{Client: mc.Client}

	bucketName, objectName := getMinioConfig(c)
	uploadID := c.GetHeader("x-amz-upload-id")

	err = coreClient.AbortMultipartUpload(context.Background(), bucketName, objectName, uploadID)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to abort multipart upload", gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Multipart upload aborted",
	})
}

// UploadFile 上传文件接口
func UploadFile(c *gin.Context) {
	mc, err := getMinioClient(c)
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
	info, err := mc.Client.PutObject(context.Background(), bucketName, objectName, src, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to upload file", gin.H{"error": err.Error()})
		return
	}

	common.RespOK(c, "File uploaded successfully", gin.H{"info": info})
}

// DeleteFile 删除文件接口
func DeleteFile(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}

	bucketName, objectName := getMinioConfig(c)

	err = mc.Client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to delete file", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "File deleted successfully", gin.H{})
}

// ListBuckets 查询bucket接口
func ListBuckets(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}

	buckets, err := mc.Client.ListBuckets(context.Background())
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to list buckets", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "List buckets successfully", gin.H{"buckets": buckets})
}

func ListObjects(c *gin.Context) {
	mc, err := getMinioClient(c)
	if err != nil {
		common.RespErr(c, http.StatusBadRequest, "Failed to create Minio client", gin.H{"error": err.Error()})
		return
	}
	bucketName, objectName := getMinioConfig(c)
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
			common.RespErr(c, http.StatusInternalServerError, "Failed to list objects", gin.H{"error": object.Err.Error()})
			return
		}
		files = append(files, object.Key)
	}
	common.RespOK(c, "List objects successfully", gin.H{"files": files})
}

// CreateNewOSSResources 创建新的OSS资源
func CreateNewOSSResourcesEndpoint(c *gin.Context) {
	var req entities.CreateNewResourcesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, http.StatusBadRequest, "Invalid request", gin.H{"error": err.Error()})
		return
	}
	userID := uint64(c.GetUint(common.UIDKey))
	if userID == 0 {
		// 使用 common.RespErr 返回错误响应
		common.RespErr(c, http.StatusBadRequest, "Failed to convert UserID to uint64", gin.H{"error": "uid is required"})
		return
	}
	// valid
	if req.Name == "" {
		common.RespErr(c, http.StatusBadRequest, "UID and Name are required", gin.H{"error": "UID and Name are required"})
		return
	}

	uid := utils.GenerateUID()
	resource, err := CreateOSS(userID, req, uid)
	if err != nil {
		common.RespErr(c, http.StatusInternalServerError, "Failed to create OSS resource", gin.H{"error": err.Error()})
		return
	}
	common.RespOK(c, "success", entities.CreateNewResourcesResponse{
		UID:  resource.UID,
		Name: resource.Name,
		Type: "oss",
	})
}

func CreateOSS(userID uint64, req entities.CreateNewResourcesRequest, uid string) (*models.OSS, error) {
	mc, err := NewMinioClient(
		fmt.Sprintf("%s:%d",
			conf.AppConfigInstance.ServerMinioHost,
			conf.AppConfigInstance.ServerMinioPort),
		conf.AppConfigInstance.ServerMinioAccess,
		conf.AppConfigInstance.ServerMinioSecret,
		conf.AppConfigInstance.ServerMinioUseSSL,
		conf.AppConfigInstance.ServerMinioRegion)
	if err != nil {
		return nil, err
	}

	resource := models.OSS{
		Name:      req.Name,
		UserID:    userID,
		UID:       uid,
		Bucket:    "",
		Region:    conf.AppConfigInstance.ServerMinioRegion,
		AccessKey: "",
		SecretKey: "",
	}
	resource.Bucket = resource.UID

	if !conf.AppConfigInstance.MinioSingleBucketMode {
		err = mc.Client.MakeBucket(context.Background(), resource.Bucket, minio.MakeBucketOptions{Region: conf.AppConfigInstance.ServerMinioRegion})
		if err != nil {
			return nil, err
		}
		account, err := CreateNewServiceAccount(resource.Bucket)
		if err != nil {
			return nil, err
		}
		resource.AccessKey = account.AccessKey
		resource.SecretKey = account.SecretKey
		resource.Expiration = account.Expiration
	} else {
		resource.SingleBucket = true
	}

	// 创建新的 OSS 资源
	db := database.GetDB()
	if err := db.Create(&resource).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

// 删除指定OSS资源
func DeleteOSSResourcesEndpoint(c *gin.Context) {
	var req entities.DeleteResourcesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, http.StatusBadRequest, "Invalid request", gin.H{"error": err.Error()})
		return
	}
	// valid
	if req.UID == "" {
		common.RespErr(c, http.StatusBadRequest, "UID is required", gin.H{"error": "UID is required"})
		return
	}

	uid := uint64(c.GetUint(common.UIDKey))
	if uid == 0 {
		common.RespErr(c, http.StatusBadRequest, "Invalid UID", gin.H{"error": "Invalid UID"})
		return
	}
	condition := models.OSS{UID: req.UID, UserID: uid}

	db := database.GetDB()
	var resource models.OSS
	// 删除OSS资源
	if rr := db.Delete(&resource, condition); rr.Error != nil || rr.RowsAffected == 0 {
		common.RespOK(c, "success but not delete bucket", entities.DeleteResourcesResp{
			Status: 0,
		})
		return
	}
	if !conf.AppConfigInstance.MinioSingleBucketMode {
		DeleteServiceAccount(resource.AccessKey)
	}

	common.RespOK(c, "success", entities.DeleteResourcesResp{
		Status: 0,
	})
}

func RecoverOSS(userID uint64, oss *models.OSS) error {
	oss.UserID = userID
	db := database.GetDB()
	oss2 := models.OSS{}
	// 如果有，则更新，如果无，则调用新增接口
	if err := db.Where("uid = ?", oss.UID).First(&models.OSS{}).Error; err != nil {
		oss.SecretKey = ""
		oss.AccessKey = ""
		oss.SessionKey = ""
		_, err := CreateOSS(userID, entities.CreateNewResourcesRequest{Name: oss.Name}, oss.UID)
		if err != nil {
			return err
		}

	} else {
		oss.SecretKey = ""
		oss.AccessKey = ""
		oss.SessionKey = ""
		if err := db.Where("uid = ?", oss.UID).Assign(oss).FirstOrCreate(&oss2).Error; err != nil {
			return err
		}
	}
	return nil
}

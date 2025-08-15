package oss

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"vvorker/conf"
	"vvorker/ext/oss/src/alioss1"
	"vvorker/funcs"

	"github.com/minio/minio-go/v7"
)

func init() {
	funcs.SetUploadFileToSysBucket(UploadFileToSysBucket)
	funcs.SetDownloadFileFromSysBucket(DownloadFileFromSysBucket)
}

func UploadFileToSysBucket(path string, obj io.Reader) error {
	switch conf.AppConfigInstance.ServerOSSType {
	case "aliyun":
		return UploadFile(path, obj)
	case "aliyun1":
		return alioss1.UploadFile(path, obj)
	default:
		return UploadFile(path, obj)
	}
}

func DownloadFileFromSysBucket(path string) (io.ReadCloser, error) {
	switch conf.AppConfigInstance.ServerOSSType {
	case "aliyun":
		return DownloadFile(path)
	case "aliyun1":
		return alioss1.DownloadFile(path)
	default:
		return DownloadFile(path)
	}
}

func getSysMinioClient() (*MinioClient, error) {
	endpoint := fmt.Sprintf("%s:%d", conf.AppConfigInstance.ServerMinioHost, conf.AppConfigInstance.ServerMinioPort)
	accessKeyID := conf.AppConfigInstance.ServerMinioAccess
	secretAccessKey := conf.AppConfigInstance.ServerMinioSecret
	region := conf.AppConfigInstance.ServerMinioRegion
	useSSLStr := fmt.Sprintf("%v", conf.AppConfigInstance.ServerMinioUseSSL)

	// logrus.Infof("Endpoint: %s, AccessKeyID: %s, SecretAccessKey: %s, UseSSL: %s, Region: %s",
	// 	endpoint, accessKeyID, secretAccessKey, useSSLStr, region)

	useSSL, err := strconv.ParseBool(useSSLStr)
	if err != nil {
		return nil, err
	}

	return NewMinioClient(endpoint, accessKeyID, secretAccessKey, useSSL, region)
}

// DownloadFileEndpoint 下载文件接口
func DownloadFile(path string) (io.ReadCloser, error) {
	mc, err := getSysMinioClient()
	if err != nil {
		return nil, err
	}

	// 使用 GetObject 获取文件流
	obj, err := mc.Client.GetObject(context.Background(),
		conf.AppConfigInstance.FileStorageOSSBucket, conf.AppConfigInstance.FileStorageOSSPrefix+path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func UploadFile(path string, obj io.Reader) error {
	mc, err := getSysMinioClient()
	if err != nil {
		return err
	}

	// 上传文件到 MinIO
	_, err = mc.Client.PutObject(context.Background(),
		conf.AppConfigInstance.FileStorageOSSBucket,
		conf.AppConfigInstance.FileStorageOSSPrefix+path,
		obj,
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return nil
}

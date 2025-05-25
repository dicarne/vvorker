package oss

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vorker/conf"
	"vorker/utils"

	"github.com/minio/madmin-go/v4"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type ServiceAccount struct {
	AccessKey    string
	SecretKey    string
	SessionToken string
	Expiration   time.Time
}

func CreateNewServiceAccount(bucket string) (ServiceAccount, error) {
	// Use a secure connection.

	creds := credentials.NewStaticV4(
		conf.AppConfigInstance.ServerMinioAccess,
		conf.AppConfigInstance.ServerMinioSecret,
		"",
	)

	mdmClnt, err := madmin.NewWithOptions(conf.AppConfigInstance.ServerMinioHost, &madmin.Options{
		Creds:  creds,
		Secure: conf.AppConfigInstance.ServerMinioUseSSL,
	})
	if err != nil {
		logrus.Errorf("Failed to create madmin client: %v", err)
		return ServiceAccount{}, err
	}

	// 定义策略
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"s3:GetBucketLocation",
					"s3:ListBucket",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s", bucket),
				},
			},
			{
				"Effect": "Allow",
				"Action": []string{
					"s3:GetObject",
					"s3:PutObject",
					"s3:DeleteObject",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", bucket),
				},
			},
		},
	}

	// 将策略转换为 JSON 字符串
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		logrus.Errorf("Failed to marshal policy: %v", err)
		return ServiceAccount{}, err
	}

	r := ServiceAccount{
		AccessKey: truncateName(utils.GenerateUID(), 20),
		SecretKey: truncateName(utils.GenerateUID(), 20),
	}

	cre, err := mdmClnt.AddServiceAccount(context.Background(), madmin.AddServiceAccountReq{
		AccessKey: r.AccessKey,
		SecretKey: r.SecretKey,
		Policy:    policyJSON,
		// 确保 Name 长度不超过 32 字符
		Name: truncateName("vorker-bucket-"+bucket, 32),
	})

	if err != nil {
		logrus.Errorf("Failed to create service account: %v", err)
		return ServiceAccount{}, err
	}
	r.SessionToken = cre.SessionToken
	r.Expiration = cre.Expiration

	return r, nil
}

// truncateName 函数用于截取字符串到指定长度
func truncateName(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

func DeleteServiceAccount(accessKey string) error {
	creds := credentials.NewStaticV4(
		conf.AppConfigInstance.ServerMinioAccess,
		conf.AppConfigInstance.ServerMinioSecret,
		"",
	)
	// Use a secure connection.
	mdmClnt, err := madmin.NewWithOptions(conf.AppConfigInstance.ServerMinioHost, &madmin.Options{
		Creds:  creds,
		Secure: conf.AppConfigInstance.ServerMinioUseSSL,
	})
	if err != nil {
		logrus.Errorf("Failed to create madmin client: %v", err)
		return err
	}
	err = mdmClnt.DeleteServiceAccount(context.Background(), accessKey)
	if err != nil {
		logrus.Errorf("Failed to delete service account: %v", err)
		return err
	}
	return nil
}

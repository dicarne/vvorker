package files

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/dao"
	oss "vvorker/ext/oss/src"
	"vvorker/models"
	"vvorker/utils"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	tarMimeType = ".tar"
	zipMimeType = ".zip"
)

type UploadFileReq struct {
	File string `json:"file" binding:"required"`
	Path string `json:"path" binding:"required"`
}

type UploadFileResp struct {
	FileID   string `json:"fileId"`
	FileHash string `json:"fileHash"`
}

func UploadFileEndpoint(c *gin.Context) {
	var req UploadFileReq
	if err := c.BindJSON(&req); err != nil {
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.File)
	if err != nil {
		logrus.WithError(err).Error("decode base64 error")
		common.RespErr(c, common.RespCodeInternalError, "Internal error processing file.", nil)
		return
	}
	contentType := filepath.Ext(req.Path)
	if contentType == zipMimeType {
		zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			logrus.WithError(err).Error("create zip reader error")
			common.RespErr(c, common.RespCodeInternalError, "Incomplete .zip archive file.", nil)
			return
		}

		data, err = utils.CreateTarFromZip(zipReader)
		if err != nil {
			logrus.WithError(err).Error("create tar from zip error")
			common.RespErr(c, common.RespCodeInternalError, "Internal error processing .zip archive.", nil)
			return
		}
		contentType = tarMimeType
	}

	hashBytes := sha256.Sum256(data)
	hash := hex.EncodeToString(hashBytes[:])
	uid := c.GetUint(common.UIDKey)
	fileRecord, err := dao.GetFileByHashAndCreator(c, hash, uid)
	if err == nil {
		if conf.AppConfigInstance.MAN_ASSET_FILE_REPLACE {
			database.GetDB().Unscoped().Model(&models.File{}).Delete(fileRecord)
		} else {
			logrus.Infof("file already exists: %s", fileRecord.UID)
			common.RespOK(c, "File already exists.", UploadFileResp{
				FileID:   fileRecord.UID,
				FileHash: hash,
			})
			return
		}

	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("get file error")
		common.RespErr(c, common.RespCodeInternalError, "Internal error getting file.", nil)
		return
	}

	if conf.AppConfigInstance.FileStorageUseOSS {
		err = oss.UploadFileToSysBucket(fmt.Sprintf("files/%d/%s", uid, hash), bytes.NewReader(data))
		if err != nil {
			logrus.WithError(err).Error("upload file to oss error")
			common.RespErr(c, common.RespCodeInternalError, "Internal error uploading file to OSS.", nil)
			return
		}
		data = nil
	}

	fileRecord, err = dao.SaveFile(c, &models.File{
		CreatedBy: uid,
		Hash:      hash,
		Mimetype:  contentType,
		Data:      data,
	})
	if err != nil {
		logrus.WithError(err).Error("insert file error")
		common.RespErr(c, common.RespCodeInternalError, "Internal error saving file.", nil)
		return
	}

	common.RespOK(c, "File uploaded successfully.", UploadFileResp{
		FileID:   fileRecord.UID,
		FileHash: hash,
	})
}

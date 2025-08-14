package files

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"path/filepath"
	"vvorker/common"
	"vvorker/dao"
	"vvorker/models"
	"vvorker/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	tarMimeType = ".tar"
	zipMimeType = ".zip"
)

type UploadFileReq struct {
	File string `json:"file"`
	Path string `json:"path"`
}

type UploadFileResp struct {
	FileID   string `json:"fileId"`
	FileHash string `json:"fileHash"`
}

func UploadFileEndpoint(c *gin.Context) {
	var req UploadFileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
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
		logrus.Infof("file already exists: %s", fileRecord.UID)
		common.RespOK(c, "File already exists.", UploadFileResp{
			FileID:   fileRecord.UID,
			FileHash: hash,
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("get file error")
		common.RespErr(c, common.RespCodeInternalError, "Internal error getting file.", nil)
		return
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

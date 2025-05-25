package files

import (
	"vvorker/common"
	"vvorker/dao"

	"github.com/gin-gonic/gin"
)

func GetFileEndpoint(c *gin.Context) {
	fileId := c.Param("fileId")
	if len(fileId) == 0 {
		common.RespErr(c, common.RespCodeInvalidRequest, "fileId is empty", nil)
		return
	}

	file, err := dao.GetFileByUID(c, c.GetUint(common.UIDKey), fileId)
	if err != nil || file == nil {
		common.RespErr(c, common.RespCodeInvalidRequest, "file not found", nil)
		return
	}

	common.RespOK(c, "get file success", file)
}

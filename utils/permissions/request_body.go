package permissions

import (
	"vvorker/common"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](c *gin.Context, req T) error {
	if err := c.BindJSON(req); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return err
	}
	return nil
}

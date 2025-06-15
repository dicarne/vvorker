package gentype

import (
	"fmt"
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/entities"
	"vvorker/ext"
	"vvorker/models"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Project struct {
	UID string `json:"uid" binding:"required"`
}

type GenTypeRequest struct {
	*conf.WorkerConfig
	Project *Project `json:"project" binding:"required"`
}

func GenerateTypes(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	worker := &GenTypeRequest{}

	if err := c.BindJSON(worker); err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}

	userID := uint64(c.GetUint(common.UIDKey))
	uid := worker.Project.UID

	logrus.Infof("userID: %d, uid: %s", userID, uid)

	db := database.GetDB()
	project := &models.Worker{}
	if err := db.Where(&models.Worker{
		Worker: &entities.Worker{
			UserID: userID,
			UID:    uid,
		},
	}).First(project).Error; err != nil {
		common.RespErr(c, common.RespCodeInvalidRequest, err.Error(), nil)
		return
	}
	finalStr := ""
	typeStr := `
export interface EnvBinding {
`
	if len(worker.WorkerConfig.Ai) > 0 {
		finalStr += ext.TypeBindingAI + "\n"
	}
	for _, v := range worker.WorkerConfig.Ai {
		typeStr += fmt.Sprintf(`
	%s: AIBinding
`, v.Binding)
	}

	if len(worker.WorkerConfig.Assets) > 0 {
		finalStr += ext.TypeBindingAssets + "\n"
	}
	for _, v := range worker.WorkerConfig.Assets {
		typeStr += fmt.Sprintf(`
	%s: AssetsBinding
`, v.Binding)
	}

	if len(worker.WorkerConfig.KV) > 0 {
		finalStr += ext.TypeBindingKV + "\n"
	}
	for _, v := range worker.WorkerConfig.KV {
		typeStr += fmt.Sprintf(`
	%s: KVBinding
`, v.Binding)
	}

	if len(worker.WorkerConfig.OSS) > 0 {
		finalStr += ext.TypeBindingOSS + "\n"
	}
	for _, v := range worker.WorkerConfig.OSS {
		typeStr += fmt.Sprintf(`
	%s: OSSBinding
`, v.Binding)
	}

	// pgsql
	if len(worker.WorkerConfig.PgSql) > 0 {
		finalStr += ext.TypeBindingPgsql + "\n"
	}
	for _, v := range worker.WorkerConfig.PgSql {
		typeStr += fmt.Sprintf(`
	%s: PGSQLBinding
`, v.Binding)
	}

	//task
	if len(worker.WorkerConfig.Task) > 0 {
		finalStr += ext.TypeBindingTask + "\n"
	}
	for _, v := range worker.WorkerConfig.Task {
		typeStr += fmt.Sprintf(`
	%s: TaskBinding
`, v.Binding)
	}

	if len(worker.WorkerConfig.Services) > 0 {
		finalStr += `
export interface AService {
		fetch: (url: string, init: RequestInit) => Promise<Response>;
}
`
		for _, v := range worker.WorkerConfig.Services {
			typeStr += fmt.Sprintf(`
	%s: AService
`, common.ToCamelCase(v))
		}
	}

	if len(worker.Vars) > 0 {
		typeStr += `
	vars: any
`
	}

	typeStr += "\n}\n" + finalStr

	common.RespOK(c, common.RespMsgOK, gin.H{
		"types": typeStr,
	})
	// c.Data(http.StatusOK, "text/plain", []byte(typeStr))
}

package export

import (
	"runtime/debug"
	"vvorker/common"
	"vvorker/conf"
	kv "vvorker/ext/kv/src"
	oss "vvorker/ext/oss/src"
	pgsql "vvorker/ext/pgsql/src"
	"vvorker/models"
	"vvorker/services/workerd"
	"vvorker/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ExportConfigReq struct {
	ServiceUIDs  []string `json:"service_uids"`  // 服务ID列表
	ServiceNames []string `json:"service_names"` // 服务名称列表
}

type AssetFile struct {
	*models.Assets
	Content []byte `json:"content"`
}

type ExportConfig struct {
	Workers []*models.Worker     `json:"workers"`
	Kv      []*models.KV         `json:"kv"`
	Pgsql   []*models.PostgreSQL `json:"pgsql"`
	Oss     []*models.OSS        `json:"oss"`
	Assets  []*AssetFile         `json:"assets"`
}

// 用于导出某些服务及所有相关资源
func ExportResourcesConfigEndpoint(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()
	var req ExportConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespErr(c, 400, "参数解析失败", nil)
		return
	}
	userID := c.GetUint(common.UIDKey)
	workers, err := models.GetWorkersByUIDs(userID, req.ServiceUIDs)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, "通过id获取worker失败", nil)
		return
	}

	workers2, err := models.GetWorkersByNames(userID, req.ServiceNames)
	if err != nil {
		common.RespErr(c, common.RespCodeInternalError, "通过name获取worker失败", nil)
		return
	}

	workers = append(workers, workers2...)

	res := ExportConfig{
		Workers: workers,
	}

	workersNameMap := make(map[string]*models.Worker)
	lastNamesQueue := make([]string, 0)
	currentNamesQueue := make([]string, 0)
	for _, w := range workers {
		workersNameMap[w.Name] = w
		ww, err := conf.ParseWorkerConfig(w.Template)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
		// 将ww.Services的name都加入lastNamesQueue
		lastNamesQueue = append(lastNamesQueue, ww.Services...)
	}

	for len(lastNamesQueue) != 0 {
		for _, name := range lastNamesQueue {
			// 如果name 不在map中
			if _, ok := workersNameMap[name]; ok {
				continue
			}
			ww, err := models.GetWorkersByNames(userID, []string{name})
			if err != nil {
				common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
				return
			}
			wc, err := conf.ParseWorkerConfig(ww[0].Template)
			if err != nil {
				common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
				return
			}
			// 将wc.Services的name都加入currentNamesQueue
			currentNamesQueue = append(currentNamesQueue, wc.Services...)
			lastNamesQueue = currentNamesQueue
			workersNameMap[name] = ww[0]
			workers = append(workers, ww[0])
		}
	}

	res.Workers = workers

	db := database.GetDB()

	for _, w := range workers {
		wc, err := conf.ParseWorkerConfig(w.Template)
		if err != nil {
			common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
		if len(wc.KV) > 0 {
			for _, ext := range wc.KV {
				if len(ext.ResourceID) != 0 {
					kvModel := &models.KV{}
					if err := db.Where(&models.KV{
						UID: ext.ResourceID,
					}).First(&kvModel).Error; err != nil {
						common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
						return
					}
					kvModel.Password = ""
					res.Kv = append(res.Kv, kvModel)
				}
			}
		}
		if len(wc.OSS) > 0 {
			for _, ext := range wc.OSS {
				if len(ext.ResourceID) != 0 {
					ossModel := &models.OSS{}
					if err := db.Where(&models.OSS{
						UID: ext.ResourceID,
					}).First(&ossModel).Error; err != nil {
						common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
						return
					}
					ossModel.SecretKey = ""
					ossModel.AccessKey = ""
					res.Oss = append(res.Oss, ossModel)
				}
			}
		}

		if len(wc.PgSql) > 0 {
			for _, ext := range wc.PgSql {
				if len(ext.ResourceID) != 0 {
					pgsqlModel := &models.PostgreSQL{}
					if err := db.Where(&models.PostgreSQL{
						UID: ext.ResourceID,
					}).First(&pgsqlModel).Error; err != nil {
						common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
						return
					}
					pgsqlModel.Password = ""
					pgsqlModel.Username = ""
					res.Pgsql = append(res.Pgsql, pgsqlModel)
				}
			}
		}

		if len(wc.Assets) > 0 {
			exportAssets := make([]*models.Assets, 0)
			if err := db.Model(&models.Assets{}).Where(&models.Assets{WorkerUID: w.UID}).Find(exportAssets).Error; err != nil {
				common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
				return
			}
			for _, asset := range exportAssets {
				file := &models.File{}
				if err := db.Where(&models.File{
					UID: asset.UID,
				}).First(file).Error; err != nil {
					common.RespErr(c, common.RespCodeInternalError, common.RespMsgInternalError, nil)
					return
				}
				assetFile := &AssetFile{
					Assets:  asset,
					Content: file.Data,
				}
				res.Assets = append(res.Assets, assetFile)
			}
		}
	}

	common.RespOK(c, "success", res)
}

func ImportResourcesConfigEndpoint(g *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered in f: %+v, stack: %+v", r, string(debug.Stack()))
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, nil)
		}
	}()

	var req ExportConfig
	if err := g.ShouldBindJSON(&req); err != nil {
		common.RespErr(g, 400, "参数解析失败", nil)
		return
	}
	userID := g.GetUint(common.UIDKey)

	for _, w := range req.Workers {
		err := workerd.Recover(userID, w.Worker)
		if err != nil {
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, err.Error())
			return
		}
	}

	for _, kv2 := range req.Kv {
		kv2.UserID = uint64(userID)
		err := kv.RecoverKV(uint64(userID), kv2)
		if err != nil {
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
	}

	for _, oss2 := range req.Oss {
		oss2.UserID = uint64(userID)
		err := oss.RecoverOSS(uint64(userID), oss2)
		if err != nil {
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
	}

	for _, pgsql2 := range req.Pgsql {
		pgsql2.UserID = uint64(userID)
		_, err := pgsql.RecoverPGSQL(uint64(userID), pgsql2)
		if err != nil {
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
	}

	db := database.GetDB()

	for _, asset := range req.Assets {
		nass := models.Assets{}
		if err := db.Where(&models.Assets{
			UID: asset.UID,
		}).Assign(asset.Assets).FirstOrCreate(&nass).Error; err != nil {
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
		nfile := models.File{}
		if err := db.Where(&models.File{
			UID: asset.UID,
		}).Assign(&models.File{
			Data:      asset.Content,
			UID:       asset.UID,
			Hash:      asset.Hash,
			Mimetype:  asset.MIME,
			CreatedBy: userID,
		}).FirstOrCreate(&nfile).Error; err != nil {
			common.RespErr(g, common.RespCodeInternalError, common.RespMsgInternalError, nil)
			return
		}
	}

	common.RespOK(g, "success", nil)
}

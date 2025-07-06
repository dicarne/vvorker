package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vvorker/conf"
	"vvorker/defs"
	"vvorker/entities"
	"vvorker/ext"
	"vvorker/funcs"
	"vvorker/utils"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

type GenTemplateConfig struct {
	Worker         *entities.Worker
	BindingsText   template.HTML
	ExtensionsText template.HTML
	ServiceText    template.HTML
	FlagsText      template.HTML
	SocketText     template.HTML
	WorkerHost     string
}

// 检查文件是否存在，若不存在则写入内容
func writeFileIfNotExists(filePath string, content string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 获取文件所在的目录
		dir := filepath.Dir(filePath)
		// 检查目录是否存在，不存在则创建
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				logrus.Errorf("Failed to create directory %s: %v", dir, err)
				return err
			}
		}
		err := utils.WriteFile(filePath, content)
		if err != nil {
			logrus.Errorf("Failed to write file %s: %v", filePath, err)
			return err
		}
	}
	return nil
}

func RPCWrapper() *req.Request {
	return req.C().R().
		SetHeaders(map[string]string{
			defs.HeaderNodeName:   conf.AppConfigInstance.NodeName,
			defs.HeaderNodeSecret: conf.RPCToken,
		})
}

func FillWorkerConfig(endpoint string, UID string) (string, error) {
	url := endpoint + "/api/agent/fill-worker-config"

	rtype := struct {
		Code int                          `json:"code"`
		Msg  string                       `json:"msg"`
		Data entities.AgentFillWorkerResp `json:"data"`
	}{}

	reqResp, err := RPCWrapper().
		SetBody(&entities.AgentFillWorkerReq{
			UID: UID,
		}).
		SetSuccessResult(&rtype).
		Post(url)

	if err != nil || reqResp.StatusCode >= 299 {
		return "", errors.New(reqResp.Err.Error())
	}
	return rtype.Data.NewTemplate, nil
}

func BuildCapfile(workers []*entities.Worker, workerQuery funcs.WorkerQuery) map[string]string {
	if len(workers) == 0 {
		return map[string]string{}
	}

	results := map[string]string{}
	for _, worker := range workers {
		writer := new(bytes.Buffer)
		capTemplate := template.New("capfile")
		workerTemplate := defs.DefaultTemplate
		newTemplate, ferr := FillWorkerConfig(conf.AppConfigInstance.MasterEndpoint, worker.GetUID())
		if ferr != nil {
			logrus.Warnf("new workerconfig error: %v", ferr)
		}
		workerconfig, werr := conf.ParseWorkerConfig(newTemplate)

		bindingsText := ""
		extensionsText := ""
		servicesText := ""
		flagsText := ""
		if werr != nil {
			logrus.Warnf("workerconfig error: %v", werr)
			workerconfig = conf.DefaultWorkerConfig()
		}

		if len(workerconfig.Extensions) > 0 {
			for _, ext := range workerconfig.Extensions {
				extName := ext.Name
				allowExtensionFn, ok := defs.AllowServicesMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					allowExtension := allowExtensionFn(ext.Binding)
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate
					switch allowExtension.Type {
					case "extension":
						extensionsText = extensionsText + ".e" + allowExtension.Name + ","
					case "worker":
						servicesText = servicesText + allowExtension.ServiceInjectTemplate
					}

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %s not found", ext)
				}
			}
		}

		if len(workerconfig.Ai) > 0 {
			for _, ext := range workerconfig.Ai {
				extName := "ai"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`( name = "API_KEY", text = "`+ext.ApiKey+`" ), ( name = "BASE_URL", text = "`+ext.BaseUrl+`" ), ( name = "MODEL", text = "`+ext.Model+`" ),`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.PgSql) > 0 {
			for _, ext := range workerconfig.PgSql {
				extName := "pgsql"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					if len(ext.ResourceID) != 0 {
						ext.Host = "localhost"
						ext.Port = conf.AppConfigInstance.ClientPostgrePort
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`
	( name = "HOST", text = "`+ext.Host+`" ), 
	( name = "PORT", text = "`+strconv.Itoa(ext.Port)+`" ), 
	( name = "USER", text = "`+ext.User+`" ),
	( name = "PASSWORD", text = "`+ext.Password+`" ),
	( name = "DATABASE", text = "`+ext.Database+`" ),`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)

				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.Mysql) > 0 {
			for _, ext := range workerconfig.Mysql {
				extName := "mysql"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					if len(ext.ResourceID) != 0 {
						ext.Host = "localhost"
						ext.Port = conf.AppConfigInstance.ClientMySQLPort
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`
	( name = "HOST", text = "`+ext.Host+`" ), 
	( name = "PORT", text = "`+strconv.Itoa(ext.Port)+`" ), 
	( name = "USER", text = "`+ext.User+`" ),
	( name = "PASSWORD", text = "`+ext.Password+`" ),
	( name = "DATABASE", text = "`+ext.Database+`" ),`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)

				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.KV) > 0 {
			for _, ext := range workerconfig.KV {
				extName := "kv"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}

					if len(ext.ResourceID) != 0 {
						ext.Host = "localhost"
						ext.Port = conf.AppConfigInstance.ClientRedisPort
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`
	( name = "HOST", text = "`+ext.Host+`" ),
	( name = "PORT", text = "`+strconv.Itoa(ext.Port)+`" ),	
	( name = "RESOURCE_ID", text = "`+ext.ResourceID+`" ),
	( name = "KVPROVIDER", text = "`+ext.Provider+`" ),
	( name = "MASTER_ENDPOINT", text = "`+conf.AppConfigInstance.MasterEndpoint+`" ),
	( name = "X_SECRET", text = "`+conf.RPCToken+`" ),
	( name = "X_NODENAME", text = "`+conf.AppConfigInstance.NodeName+`" ),
`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.OSS) > 0 {
			for _, ext := range workerconfig.OSS {
				extName := "oss"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					ossAgentUrl := conf.AppConfigInstance.MasterEndpoint
					if len(ext.ResourceID) != 0 {
						ossAgentUrl = fmt.Sprintf("http://127.0.0.1:%d", conf.AppConfigInstance.APIPort)
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`
	( name = "HOST", text = "`+ext.Host+`" ),
	( name = "PORT", text = "`+strconv.Itoa(ext.Port)+`" ),
	( name = "ACCESS_KEY_ID", text = "`+ext.AccessKeyId+`" ),
	( name = "ACCESS_KEY_SECRET", text = "`+ext.AccessKeySecret+`" ),
	( name = "BUCKET", text = "`+ext.Bucket+`" ),
	( name = "USE_SSL", text = "`+strconv.FormatBool(ext.UseSSL)+`" ),
	( name = "REGION", text = "`+ext.Region+`" ),
	( name = "OSS_AGENT_URL", text = "`+ossAgentUrl+`" ),
	( name = "RESOURCE_ID", text = "`+ext.ResourceID+`" ),
	( name = "X_SECRET" , text = "`+conf.RPCToken+`" ),
	( name = "X_NODENAME", text = "`+conf.AppConfigInstance.NodeName+`" ),
`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.Assets) > 0 {
			for _, ext := range workerconfig.Assets {
				extName := "assets"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`
	( name = "WORKER_UID", text = "`+worker.UID+`" ),
	( name = "MASTER_ENDPOINT", text = "`+conf.AppConfigInstance.MasterEndpoint+`" ),
	( name = "X_SECRET" , text = "`+conf.RPCToken+`" ),
	( name = "X_NODENAME", text = "`+conf.AppConfigInstance.NodeName+`" ),
`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.Task) > 0 {
			for _, ext := range workerconfig.Task {
				extName := "task"
				allowExtensionFn, ok := defs.AllowWorkersMap[extName]
				if ok {
					if len(ext.Binding) == 0 {
						ext.Binding = extName
					}
					allowExtension := allowExtensionFn(ext.Binding, template.HTML(`
	( name = "WORKER_UID", text = "`+worker.UID+`" ),
	( name = "MASTER_ENDPOINT", text = "`+conf.AppConfigInstance.MasterEndpoint+`" ),
	( name = "X_SECRET" , text = "`+conf.RPCToken+`" ),
	( name = "X_NODENAME", text = "`+conf.AppConfigInstance.NodeName+`" ),
`))
					workerTemplate = workerTemplate + allowExtension.ExtensionTemplate
					bindingsText = bindingsText + allowExtension.BindingTemplate

					servicesText = servicesText + allowExtension.ServiceInjectTemplate

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "../../lib", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %v not found", ext)
				}
			}
		}

		if len(workerconfig.Services) > 0 {
			for _, service := range workerconfig.Services {
				netw := defs.GenServiceNetwork(worker.UID, service, workerQuery)
				workerTemplate = workerTemplate + netw.NetworkText
				servicesText = servicesText + netw.ServiceText
				bindingsText = bindingsText + netw.BindingsText
			}
		}

		if len(workerconfig.CompatibilityFlags) > 0 {
			for _, flag := range workerconfig.CompatibilityFlags {
				flagsText = flagsText + "\"" + flag + "\","
			}
		}

		if len(workerconfig.Vars) > 0 {
			jsonBytes, err := json.Marshal(string(workerconfig.Vars))
			if err != nil {
				logrus.Errorln("Error:", err)
			} else {
				jsonString := string(jsonBytes)
				bindingsText += "( name = \"vars\", json = " + jsonString + " ),"
			}
		}

		servicesText += defs.DefaultControlService
		controlWorker := strings.ReplaceAll(defs.DefaultControlWorker, "{{.BindingsMainWorker}}", `(name = "worker", service = "`+worker.UID+`"),`)
		workerTemplate += controlWorker

		writeFileIfNotExists(filepath.Join(conf.AppConfigInstance.WorkerdDir,
			defs.WorkerInfoPath,
			worker.GetUID(), "../../lib", "control.js"),
			ext.ControlScript)

		capTemplate, err := capTemplate.Parse(workerTemplate)
		if err != nil {
			panic(err)
		}

		socketText := defs.DefaultSocketText
		socketText = strings.ReplaceAll(socketText, "{{.ControlHost}}", "localhost")
		socketText = strings.ReplaceAll(socketText, "{{.ControlPort}}", strconv.Itoa(int(worker.ControlPort)))

		genConfig := GenTemplateConfig{
			Worker:         worker,
			BindingsText:   template.HTML(bindingsText),
			ExtensionsText: template.HTML(extensionsText),
			ServiceText:    template.HTML(servicesText),
			FlagsText:      template.HTML(flagsText),
			SocketText:     template.HTML(socketText),
			WorkerHost:     worker.Name + conf.AppConfigInstance.WorkerURLSuffix,
		}
		capTemplate.Execute(writer, genConfig)
		results[worker.GetUID()] = writer.String()
	}
	return results
}

func GenWorkerConfig(worker *entities.Worker, workerQuery funcs.WorkerQuery) error {
	if worker == nil || worker.GetUID() == "" {
		return errors.New("error worker")
	}
	fileMap := BuildCapfile([]*entities.Worker{
		worker,
	}, workerQuery)

	fileContent, ok := fileMap[worker.GetUID()]
	if !ok {
		return errors.New("BuildCapfile error")
	}

	return utils.WriteFile(
		filepath.Join(
			conf.AppConfigInstance.WorkerdDir,
			defs.WorkerInfoPath,
			worker.GetUID(),
			defs.CapFileName,
		), fileContent)
}

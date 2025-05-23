package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"vorker/conf"
	"vorker/defs"
	"vorker/entities"

	"github.com/sirupsen/logrus"
)

type GenTemplateConfig struct {
	Worker         *entities.Worker
	BindingsText   template.HTML
	ExtensionsText template.HTML
	ServiceText    template.HTML
	FlagsText      template.HTML
}

// 检查文件是否存在，若不存在则写入内容
func writeFileIfNotExists(filePath string, content string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := WriteFile(filePath, content)
		if err != nil {
			logrus.Errorf("Failed to write file %s: %v", filePath, err)
		}
	}
}

func BuildCapfile(workers []*entities.Worker) map[string]string {
	if len(workers) == 0 {
		return map[string]string{}
	}

	results := map[string]string{}
	for _, worker := range workers {
		writer := new(bytes.Buffer)
		capTemplate := template.New("capfile")
		workerTemplate := defs.DefaultTemplate
		workerconfig, werr := conf.ParseWorkerConfig(worker.Template)

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
					if allowExtension.Type == "extension" {
						extensionsText = extensionsText + ".e" + allowExtension.Name + ","
					} else if allowExtension.Type == "worker" {
						servicesText = servicesText + allowExtension.ServiceInjectTemplate
					}

					// 构建文件路径
					filePath := filepath.Join(conf.AppConfigInstance.WorkerdDir,
						defs.WorkerInfoPath,
						worker.GetUID(), "src", extName+".js")

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
						worker.GetUID(), "src", extName+".js")

					writeFileIfNotExists(filePath, allowExtension.Script)
				} else {
					logrus.Warnf("service %s not found", ext)
				}
			}
		}

		if len(workerconfig.Services) > 0 {
			for _, service := range workerconfig.Services {
				netw := defs.GenServiceNetwork(service)
				workerTemplate = workerTemplate + netw.NetworkText
				servicesText = servicesText + netw.ServiceText
				bindingsText = bindingsText + netw.BindingsText
			}
		}

		if len(workerconfig.CompatibilityFlags) > 0 {
			for _, flag := range workerconfig.CompatibilityFlags {
				flagsText = flagsText + flag + ","
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

		capTemplate, err := capTemplate.Parse(workerTemplate)
		if err != nil {
			panic(err)
		}

		genConfig := GenTemplateConfig{
			Worker:         worker,
			BindingsText:   template.HTML(bindingsText),
			ExtensionsText: template.HTML(extensionsText),
			ServiceText:    template.HTML(servicesText),
			FlagsText:      template.HTML(flagsText),
		}
		capTemplate.Execute(writer, genConfig)
		results[worker.GetUID()] = writer.String()
	}
	return results
}

func GenWorkerConfig(worker *entities.Worker) error {
	if worker == nil || worker.GetUID() == "" {
		return errors.New("error worker")
	}
	fileMap := BuildCapfile([]*entities.Worker{
		worker,
	})

	fileContent, ok := fileMap[worker.GetUID()]
	if !ok {
		return errors.New("BuildCapfile error")
	}

	return WriteFile(
		filepath.Join(
			conf.AppConfigInstance.WorkerdDir,
			defs.WorkerInfoPath,
			worker.GetUID(),
			defs.CapFileName,
		), fileContent)
}

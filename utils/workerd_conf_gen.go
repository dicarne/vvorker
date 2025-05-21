package utils

import (
	"bytes"
	"errors"
	"html/template"
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
		if werr != nil {
			logrus.Warnf("workerconfig error: %v", werr)
			workerconfig = conf.DefaultWorkerConfig()
		}

		if len(workerconfig.EnabledServices) > 0 {
			for _, service := range workerconfig.EnabledServices {
				allowService, ok := defs.AllowServicesMap[service]
				if ok {
					workerTemplate = workerTemplate + allowService.ServiceTemplate
					bindingsText = bindingsText + allowService.BindingTemplate
					extensionsText = extensionsText + "." + allowService.Name + ","
				} else {
					logrus.Warnf("service %s not found", service)
				}
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

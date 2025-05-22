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
	ServiceText    template.HTML
	FlagsText      template.HTML
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
			for _, service := range workerconfig.Extensions {
				allowService, ok := defs.AllowServicesMap[service]
				if ok {
					workerTemplate = workerTemplate + allowService.ExtensionTemplate
					bindingsText = bindingsText + allowService.BindingTemplate
					if allowService.Type == "extension" {
						extensionsText = extensionsText + "." + allowService.Name + ","
					} else if allowService.Type == "worker" {
						servicesText = servicesText + allowService.ServiceInjectTemplate
					}
				} else {
					logrus.Warnf("service %s not found", service)
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

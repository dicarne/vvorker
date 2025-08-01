package defs

import (
	"bytes"
	"html/template"
	"vvorker/ext"

	"github.com/sirupsen/logrus"
)

type AllowServiceTemplate struct {
	Name                  string
	ExtensionTemplate     string
	ServiceInjectTemplate string
	BindingTemplate       string
	Path                  string
	BasicServiceTemplate  string
	BasicBindingTemplate  string
	Type                  string
	Script                string
	WorkerBindingText     template.HTML
	FlagsText             template.HTML
}

var commonExtensionTemplate = `
const e{{.Name}} :Workerd.Extension = (
  modules = [
    (name = "e{{.Name}}:binding", esModule = embed "../../lib/{{.Path}}.js", internal = true),
  ],
);
`

var commonWorkerTemplate = `
 const w{{.Name}} :Workerd.Worker = (
   modules = [
     (name = "{{.Name}}", esModule = embed "../../lib/{{.Path}}.js"),
   ],
   compatibilityDate = "2025-05-08",
   bindings = [{{.WorkerBindingText}}],
   compatibilityFlags = [{{.FlagsText}}],
 );
 `

var commonExtensionBindingTemplate = `
	(name = "{{.Name}}", wrapped = ( moduleName = "e{{.Name}}:binding" )),
`

var commonWorkerBindingTemplate = `
	(name = "{{.Name}}", service = "{{.Name}}"),
`

var commonServiceInjectTemplate = `
	(name = "{{.Name}}", worker = .w{{.Name}}),
`

// 生成简单扩展模板
func GenerateExtensionTemplate(temp AllowServiceTemplate) AllowServiceTemplate {
	capTemplate := template.New("basic")
	basicTemplate, err := capTemplate.Parse(temp.BasicServiceTemplate)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
		panic(err)
	}
	writer1 := new(bytes.Buffer)
	basicTemplate.Execute(writer1, temp)
	temp.ExtensionTemplate = writer1.String()

	bindingTemplate := template.New("binding")
	basicBinding, err := bindingTemplate.Parse(temp.BasicBindingTemplate)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
		panic(err)
	}
	writer2 := new(bytes.Buffer)
	basicBinding.Execute(writer2, temp)
	temp.BindingTemplate = writer2.String()

	serviceInjectTemplate := template.New("serviceInject")
	basicServiceInject, err := serviceInjectTemplate.Parse(temp.ServiceInjectTemplate)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
		panic(err)
	}
	writer3 := new(bytes.Buffer)
	basicServiceInject.Execute(writer3, temp)
	temp.ServiceInjectTemplate = writer3.String()
	return temp
}

var AllowServicesMap = map[string]func(name string) AllowServiceTemplate{}

var AllowWorkersMap = map[string]func(name string, workerBindings template.HTML) AllowServiceTemplate{
	"ai": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "ai",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtAiScript,
			WorkerBindingText:     workerBindings,
		})
	},
	"pgsql": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "pgsql",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtPgsqlScript,
			FlagsText:             template.HTML(` "nodejs_compat" `),
			WorkerBindingText:     workerBindings,
		})
	},
	"oss": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "oss",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtOSSScript,
			FlagsText:             template.HTML(` "nodejs_compat" `),
			WorkerBindingText:     workerBindings,
		})
	},
	"kv": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "kv",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtKVScript,
			FlagsText:             template.HTML(` "nodejs_compat" `),
			WorkerBindingText:     workerBindings,
		})
	},
	"assets": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "assets",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtAssetsScript,
			FlagsText:             template.HTML(``),
			WorkerBindingText:     workerBindings,
		})
	},
	"task": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "task",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtTaskScript,
			FlagsText:             template.HTML(``),
			WorkerBindingText:     workerBindings,
		})
	},
	"mysql": func(name string, workerBindings template.HTML) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                  name,
			Path:                  "mysql",
			BasicServiceTemplate:  commonWorkerTemplate,
			BasicBindingTemplate:  commonWorkerBindingTemplate,
			ServiceInjectTemplate: commonServiceInjectTemplate,
			Type:                  "worker",
			Script:                ext.ExtMysqlScript,
			FlagsText:             template.HTML(` "nodejs_compat" `),
			WorkerBindingText:     workerBindings,
		})
	},
}

type ServiceNetworkTemplate struct {
	NetworkText  string
	ServiceText  string
	BindingsText string
	SocketText   string
}

type NameStruct struct {
	Name string
}

package defs

import (
	"bytes"
	"html/template"
	"vorker/ext"

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
}

var commonBasicServiceTemplate = `
const e{{.Name}} :Workerd.Extension = (
  modules = [
    (name = "e{{.Name}}:binding", esModule = embed "src/{{.Path}}.js", internal = true),
  ],
);
`

var commonBasicBindingTemplate = `
(name = "{{.Name}}", wrapped = (
	moduleName = "e{{.Name}}:binding"
	)
),
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

var AllowServicesMap = map[string]func(name string) AllowServiceTemplate{
	"ai": func(name string) AllowServiceTemplate {
		return GenerateExtensionTemplate(AllowServiceTemplate{
			Name:                 name,
			BasicServiceTemplate: commonBasicServiceTemplate,
			BasicBindingTemplate: commonBasicBindingTemplate,
			Type:                 "extension",
			Script:               ext.ExtAiScript,
			Path:                 "ai",
		})
	},
}

type ServiceNetworkTemplate struct {
	NetworkText  string
	ServiceText  string
	BindingsText string
}

type NameStruct struct {
	Name string
}

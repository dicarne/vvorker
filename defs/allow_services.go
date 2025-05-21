package defs

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"
)

type AllowServiceTemplate struct {
	Name            string
	ServiceTemplate string
	BindingTemplate string
	Path            string
}

var basicServiceTemplate = `
const {{.Name}} :Workerd.Extension = (
  modules = [
    (name = "{{.Name}}:binding", esModule = embed "{{.Path}}", internal = true),
  ],
);
`

var basicBindingTemplate = `
(name = "{{.Name}}", wrapped = (
	moduleName = "{{.Name}}:binding"
	)
),
`

func GenerateTemplate(temp AllowServiceTemplate) AllowServiceTemplate {
	capTemplate := template.New("basic")
	basicTemplate, err := capTemplate.Parse(basicServiceTemplate)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
		panic(err)
	}
	writer1 := new(bytes.Buffer)
	basicTemplate.Execute(writer1, temp)
	temp.ServiceTemplate = writer1.String()

	bindingTemplate := template.New("binding")
	basicBinding, err := bindingTemplate.Parse(basicBindingTemplate)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
		panic(err)
	}
	writer2 := new(bytes.Buffer)
	basicBinding.Execute(writer2, temp)
	temp.BindingTemplate = writer2.String()
	return temp
}

var AllowServicesMap = map[string]AllowServiceTemplate{
	"ai": GenerateTemplate(AllowServiceTemplate{Name: "ai", Path: "../../libs/openai/dist/index.js"}),
}

package defs

import (
	"bytes"
	"html/template"
	"strings"
	"vorker/conf"

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
}

var commonBasicServiceTemplate = `
const {{.Name}} :Workerd.Extension = (
  modules = [
    (name = "{{.Name}}:binding", esModule = embed "{{.Path}}", internal = true),
  ],
);
`

var commonBasicBindingTemplate = `
(name = "{{.Name}}", wrapped = (
	moduleName = "{{.Name}}:binding"
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

var AllowServicesMap = map[string]AllowServiceTemplate{
	"ai": GenerateExtensionTemplate(AllowServiceTemplate{
		Name:                  "ai",
		Path:                  "../../libs/openai/dist/index.js",
		BasicServiceTemplate:  commonBasicServiceTemplate,
		BasicBindingTemplate:  commonBasicBindingTemplate,
		ServiceInjectTemplate: "",
		Type:                  "extension",
	}),
}

type ServiceNetworkTemplate struct {
	NetworkText  string
	ServiceText  string
	BindingsText string
}

type NameStruct struct {
	Name string
}

// 用于生成内部网络绑定的代码。会生成将服务名称改为驼峰命名后的env绑定。
//
// @param serviceName 服务名称，例如：test-service
func GenServiceNetwork(serviceName string) ServiceNetworkTemplate {
	name := strings.ReplaceAll(serviceName, "-", "_")

	networkTemplate := template.New("network")
	networkTemplateW, err := networkTemplate.Parse(`
const n{{.Name}}Network :Workerd.ExternalServer = (
  address = "127.0.0.1:8080",
  http = (
	injectRequestHeaders = [
	  (name = "host", value = "{{.Host}}{{.Domain}}"),
	]
  )
);
`)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
	}
	writer := new(bytes.Buffer)
	networkTemplateW.Execute(writer, struct {
		Name   string
		Host   string
		Domain string
	}{
		Name:   name,
		Host:   serviceName,
		Domain: conf.AppConfigInstance.WorkerURLSuffix,
	})

	// ----------

	serviceTemplate := template.New("service")
	serviceTemplateW, err := serviceTemplate.Parse(`
(name = "{{.Name}}Network", external = .n{{.Name}}Network),
	`)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
	}
	writer2 := new(bytes.Buffer)
	serviceTemplateW.Execute(writer2, NameStruct{
		Name: name,
	})

	// ----------

	bindingTemplate := template.New("binding")
	bindingTemplateW, err := bindingTemplate.Parse(`
	(name = "{{.Name}}", service="{{.Name}}Network"),
	`)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
	}
	writer3 := new(bytes.Buffer)
	bindingTemplateW.Execute(writer3, NameStruct{
		Name: name,
	})

	return ServiceNetworkTemplate{
		NetworkText:  writer.String(),
		ServiceText:  writer2.String(),
		BindingsText: writer3.String(),
	}
}

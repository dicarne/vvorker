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

func GenerateTemplate(temp AllowServiceTemplate) AllowServiceTemplate {
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
	"ai": GenerateTemplate(AllowServiceTemplate{
		Name:                  "ai",
		Path:                  "../../libs/openai/dist/index.js",
		BasicServiceTemplate:  commonBasicServiceTemplate,
		BasicBindingTemplate:  commonBasicBindingTemplate,
		ServiceInjectTemplate: "",
		Type:                  "extension",
	}),
	// 	"request": GenerateTemplate(AllowServiceTemplate{
	// 		Name: "request",
	// 		Path: "../../libs/request/dist/index.js",
	// 		BasicServiceTemplate: `
	// const {{.Name}} :Workerd.Worker = (
	//   modules = [
	//     (name = "{{.Name}}", esModule = embed "{{.Path}}"),
	//   ],
	//   compatibilityDate = "2025-05-01",
	//   bindings = [
	//     (name = "internalNet", service = "vorkerInternalNetwork"),
	//   ],
	// );
	// `,
	// 		BasicBindingTemplate: `
	// (name = "{{.Name}}", service = "{{.Name}}"),
	// `,
	// 		ServiceInjectTemplate: `(name = "{{.Name}}", worker= .{{.Name}} ),`,
	// 		Type:                  "worker",
	// 	}),
}

type ServiceNetworkTemplate struct {
	NetworkText  string
	ServiceText  string
	BindingsText string
}

type NameStruct struct {
	Name string
}

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

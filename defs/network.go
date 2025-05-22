package defs

import (
	"bytes"
	"html/template"
	"strings"
	"vorker/conf"

	"github.com/sirupsen/logrus"
)

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

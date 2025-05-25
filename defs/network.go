package defs

import (
	"bytes"
	"html/template"
	"strings"
	"vvorker/conf"

	"github.com/sirupsen/logrus"
)

// 将字符串转换为驼峰命名
func toCamelCase(s string) string {
	var result string
	capitalizeNext := false
	for i, char := range s {
		if char == '-' || char == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext || i == 0 {
				result += strings.ToUpper(string(char))
				capitalizeNext = false
			} else {
				result += strings.ToLower(string(char))
			}
		}
	}
	return result
}

// 用于生成内部网络绑定的代码。会生成将服务名称改为驼峰命名后的env绑定。
//
// @param serviceName 服务名称，例如：test-service
func GenServiceNetwork(serviceName string) ServiceNetworkTemplate {
	// 修改为转换为驼峰命名
	name := toCamelCase(serviceName)

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

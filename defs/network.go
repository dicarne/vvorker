package defs

import (
	"bytes"
	"html/template"
	"vvorker/common"
	"vvorker/conf"
	"vvorker/funcs"

	"github.com/sirupsen/logrus"
)

// 用于生成内部网络绑定的代码。会生成将服务名称改为驼峰命名后的env绑定。
//
// @param serviceName 服务名称，例如：test-service
func GenServiceNetwork(thisWorkerUID string, serviceName string, w funcs.WorkerQuery) ServiceNetworkTemplate {
	// 修改为转换为驼峰命名
	name := common.ToCamelCase(serviceName)

	networkTemplate := template.New("network")
	networkTemplateW, err := networkTemplate.Parse(`
const n{{.Name}}Network :Workerd.ExternalServer = (
  address = "127.0.0.1:{{.Port}}",
  http = (
	injectRequestHeaders = [
	  (name = "host", value = "{{.Host}}{{.Domain}}"),
	  (name = "vvorker-internal-token", value = "{{.Token}}"),
	  (name = "vvorker-worker-uid", value = "{{.WorkerUID}}"),
	],
	{{.Capnp}}
  )
);
`)
	if err != nil {
		logrus.Errorf("Failed to parse template: %v", err)
	}
	port := 0
	capnp := "capnpConnectHost = \"" + serviceName + conf.AppConfigInstance.WorkerURLSuffix + "\","
	port, err = w.WorkerNameToPort(serviceName)
	if err != nil {
		port = conf.AppConfigInstance.WorkerPort
		capnp = ""
	}
	writer := new(bytes.Buffer)
	networkTemplateW.Execute(writer, struct {
		Name      string
		Host      string
		Domain    string
		Token     string
		WorkerUID string
		Port      int
		Capnp     template.HTML
	}{
		Name:      name,
		Host:      serviceName,
		Domain:    conf.AppConfigInstance.WorkerURLSuffix,
		Token:     thisWorkerUID + ":" + conf.RPCToken,
		WorkerUID: thisWorkerUID,
		Port:      port,
		Capnp:     template.HTML(capnp),
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

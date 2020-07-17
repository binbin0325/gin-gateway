//Alibaba Nacos config center
package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/assert"
	"testing"
)
var serverConfigTest = constant.ServerConfig{
	ContextPath: "/nacos",
	Port:        80,
	IpAddr:      "console.nacos.io",
}
var clientConfigTest = constant.ClientConfig{
	TimeoutMs:      10000,
	ListenInterval: 20000,
	BeatInterval:   10000,
}

func cretateConfigClientTest() config_client.ConfigClient {
	nc := nacos_client.NacosClient{}
	nc.SetServerConfig([]constant.ServerConfig{serverConfigTest})
	nc.SetClientConfig(clientConfigTest)
	nc.SetHttpAgent(&http_agent.HttpAgent{})
	client, _ := config_client.NewConfigClient(&nc)
	return client
}

func TestNacosGetConfigByNamespaceErr(t *testing.T){
	clientConfigTest.NamespaceId="test1"
	configClient := cretateConfigClientTest()
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "test_data",
		Group:  "test_group",
	})
	assert.True(t, content == "" && err != nil, "NacosDataSource get config failed.")
}

func TestNacosGetConfigByNamespace(t *testing.T){
	clientConfigTest.NamespaceId="test1"
	configClient := cretateConfigClientTest()
	published, err := configClient.PublishConfig(vo.ConfigParam{
		DataId: "test_data",
		Group:  "test_group",
		Content: "xxxxx",
	})
	assert.True(t, published && err == nil, "NacosDataSource get config failed.")
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "test_data",
		Group:  "test_group",
	})
	fmt.Println(content)
	assert.True(t, content != "" && err == nil, "NacosDataSource get config sucess.")
}
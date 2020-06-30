//Alibaba Nacos config center
package nacos

import (
	"fmt"
	"gin-gateway/pkg/cc"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var clientConfig = constant.ClientConfig{
	TimeoutMs:           10 * 1000,
	BeatInterval:        5 * 1000,
	ListenInterval:      300 * 1000,
	NotLoadCacheAtStart: true,
}

type ConfigServer struct {
	Ip       string
	Port     uint64
	Username string
	Password string
}

var once sync.Once
var configClient config_client.ConfigClient

func (cs *ConfigServer) GetConfigClient() interface{} {
	//实现单例
	once.Do(func() {
		configClient = cs.initConfigClient()
	})
	return configClient

}

func (cs *ConfigServer) initConfigClient() config_client.ConfigClient {
	nc := nacos_client.NacosClient{}
	nc.SetServerConfig([]constant.ServerConfig{constant.ServerConfig{
		IpAddr:      cs.Ip,
		Port:        cs.Port,
		ContextPath: "/nacos",
	}})
	clientConfig.Password = cs.Password
	clientConfig.Username = cs.Username
	nc.SetClientConfig(clientConfig)
	nc.SetHttpAgent(&http_agent.HttpAgent{})
	client, _ := config_client.NewConfigClient(&nc)
	return client
}

func main() {
	var zz cc.ConfigCenter
	zz = &ConfigServer{
		Ip:   viper.GetString("nacos.config.ip"),
		Port: viper.GetUint64("nacos.config.port"),
	}
	client := zz.GetConfigClient().(config_client.ConfigClient)
	content, _ := client.GetConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	})
	fmt.Println("config :" + content)
	_, err := client.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "group",
		Content: "hello world!"})
	if err != nil {
		fmt.Printf("success err:%s", err.Error())
	}
	content = ""

	client.ListenConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", data:" + data)
			content = data
		},
	})

	client.ListenConfig(vo.ConfigParam{
		DataId: "abc",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

	time.Sleep(5 * time.Second)
	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "group",
		Content: "abc"})

	select {}

}

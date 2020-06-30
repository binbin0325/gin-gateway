package nacos

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/spf13/viper"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, _ := client.RegisterInstance(param)
	fmt.Println(success)
}

func DeRegisterServiceInstance(client naming_client.INamingClient, param vo.DeregisterInstanceParam) {
	success, _ := client.DeregisterInstance(param)
	fmt.Println(success)
}

func GetService(client naming_client.INamingClient) {
	service, _ := client.GetService(vo.GetServiceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"a"},
	})
	fmt.Println(utils.ToJsonString(service))
}

func Subscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
	client.Subscribe(param)
}

func UnSubscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
	client.Unsubscribe(param)
}

func NacosRegister() {
	client := GetRegisterClient()
	var (
		port        uint64
		weight      float64
		clusterName string
		enable      bool
		healthy     bool
		ephemeral   bool
		clientIp    string
	)
	serviceName := viper.GetString("server.name")
	if viper.GetUint64("server.port") != 0 {
		port = viper.GetUint64("server.port")
	} else {
		port = 8080
	}
	if viper.GetFloat64("nacos.discovery.weight") != 0 {
		weight = viper.GetFloat64("nacos.discovery.weight")
	} else {
		weight = 10
	}
	if viper.GetString("nacos.discovery.clusterName") != "" {
		clusterName = viper.GetString("nacos.discovery.clusterName")
	} else {
		clusterName = "default"
	}
	if viper.GetBool("nacos.discovery.enable") {
		enable = viper.GetBool("nacos.discovery.enable")
	} else {
		enable = true
	}
	if viper.GetBool("nacos.discovery.healthy") {
		healthy = viper.GetBool("nacos.discovery.healthy")
	} else {
		healthy = true
	}
	if viper.GetBool("nacos.discovery.ephemeral") {
		ephemeral = viper.GetBool("nacos.discovery.ephemeral")
	} else {
		ephemeral = true
	}
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						clientIp = ipnet.IP.String()
					}
				}
			}
		}
	}
	RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          clientIp,
		Port:        port,
		ServiceName: serviceName,
		Weight:      weight,
		ClusterName: clusterName,
		Enable:      enable,
		Healthy:     healthy,
		Ephemeral:   ephemeral,
	})
}

var once sync.Once
var namingClient naming_client.INamingClient

func GetRegisterClient() naming_client.INamingClient {
	//实现单例
	once.Do(func() {
		namingClient = initRegisterClient()
	})
	return namingClient

}

func initRegisterClient() naming_client.INamingClient {

	discoveryIp := viper.GetString("nacos.discovery.ip")
	discoveryPort := viper.GetUint64("nacos.discovery.port")
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{
			{
				IpAddr: discoveryIp,
				Port:   discoveryPort,
			},
		},
		"clientConfig": constant.ClientConfig{
			TimeoutMs:           5000,
			ListenInterval:      10000,
			NotLoadCacheAtStart: true,
			LogDir:              viper.GetString("nacos.log"),
			CacheDir:             viper.GetString("nacos.cache"),
			//Username:			 "nacos",
			//Password:			 "nacos",
		},
	})

	if err != nil {
		panic(err)
	}
	return client
}

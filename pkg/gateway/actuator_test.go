package gateway

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initRegisterClient() naming_client.INamingClient {

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{
			{
				IpAddr: "192.168.23.178",
				Port:   8848,
			},
		},
		"clientConfig": constant.ClientConfig{
			TimeoutMs:           5000,
			ListenInterval:      10000,
			NotLoadCacheAtStart: true,
		},
	})

	if err != nil {
		panic(err)
	}
	return client
}

func TestGetInstance(t *testing.T) {
	services, err := initRegisterClient().GetAllServicesInfo(vo.GetAllServiceInfoParam{
		GroupName: "DEFAULT_GROUP",
	})
	fmt.Println(services)
	assert.True(t, err == nil, true)
}

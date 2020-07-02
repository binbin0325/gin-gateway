package gateway

import (
	"encoding/json"
	"gin-gateway/pkg/cc"
	nacos_config "gin-gateway/pkg/cc/nacos"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Router struct {
	Uri        string      `json:"uri"`
	Type       string      `json:"type,omitempty"`
	Predicates []Predicate `json:"predicates,omitempty"`
	Order      int         `json:"order,omitempty"`
	Filters    []Filter    `json:"filters,omitempty"`
	RouterFiltersChain RouterHandlersChain
}

type Predicate struct {
	Args map[string]string `json:"args,omitempty"`
	Name string            `json:"name,omitempty"`
}

type Filter struct {
	Name string            `json:"name,omitempty"`
	Args map[string]string `json:"args,omitempty"`
}

var contextPath string

var routerMap map[string]*Router

func InitRouter(router *gin.Engine) {
	contextPath = viper.GetString("server.router.context_path")
	routerGroup := router.Group(contextPath)
	initGlobalFilters(routerGroup)
	loadRouter(getRouters(), routerGroup)
	initRouterFiltersChain()
	ginpprof.Wrap(router)

}
func getRouters() []*Router {
	var cc cc.ConfigCenter
	cc = &nacos_config.ConfigServer{
		Ip:   viper.GetString("nacos.config.ip"),
		Port: viper.GetUint64("nacos.config.port"),
	}
	configClient := cc.GetConfigClient().(config_client.ConfigClient)
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "routers",
		Group:  "gogateway",
	})
	if err != nil {
		os.Exit(1)
	}
	var routers []*Router
	err = json.Unmarshal([]byte(content), &routers)
	if err != nil {
		os.Exit(1)
	}
	return routers
}

func loadRouter(routers []*Router, routerGroup *gin.RouterGroup) {
	routerMapping := make(map[string]*Router)
	for _, r := range routers {
		for _, p := range r.Predicates {
			pattern := p.Args["pattern"]
			if index := strings.Index(pattern, "*"); index > 0 {
				requestKey := pattern[:index]
				routerMapping[contextPath+requestKey] = r
				pattern = pattern + "path"
			} else {
				routerMapping[contextPath+pattern] = r
			}
			routerGroup.Any(pattern, actuator)
		}
	}
	routerMap = routerMapping
}



package gateway

import (
	"encoding/json"
	"fmt"
	"gin-gateway/pkg/cc"
	nacos_config "gin-gateway/pkg/cc/nacos"
	"gin-gateway/pkg/rc/nacos"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

type Router struct {
	Uri        string      `json:"uri"`
	Type       string      `json:"type,omitempty"`
	Predicates []Predicate `json:"predicates,omitempty"`
	Order      int         `json:"order,omitempty"`
	Filters    []Filter    `json:"filters,omitempty"`
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
	routerMap = loadRouter(getRouters(), routerGroup)
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

func loadRouter(routers []*Router, routerGroup *gin.RouterGroup) map[string]*Router {
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
	return routerMapping
}

func actuator(c *gin.Context) {
	for k, v := range routerMap {
		if strings.Contains(c.Request.RequestURI, k) {
			v.proxy(c)
		}
	}
}

func (v *Router) proxy(c *gin.Context) {
	var host string
	if v.Type == "lb" {
		instance := getInstance(v.Uri)
		host = instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
	}
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.URL.Path = c.Request.URL.Path
		targetQuery:=c.Request.URL.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)

}

func getInstance(uri string) *model.Instance {
	instance, err := nacos.GetRegisterClient().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: uri,
	})
	if err != nil {
		fmt.Println(err)
	}
	return instance

}

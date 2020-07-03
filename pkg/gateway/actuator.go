package gateway

import (
	"fmt"
	"gin-gateway/pkg/rc/nacos"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

func actuator(c *gin.Context) {
	for k, v := range routerMap {
		//todo strings.Contains update
		if strings.Contains(c.Request.RequestURI, k) {
			executeRouterFiltersChain(v, c)
			v.proxy(c)
		}
	}
}

func (v *Router) proxy(c *gin.Context) {
	var host string
	//todo not lb
	if v.Type == "lb" {
		if instance := getInstance(v.Uri); instance != nil {
			host = instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
		} else {
			c.JSON(http.StatusInternalServerError, "not find instance: "+v.Uri)
		}
	}
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.URL.Path = c.Request.URL.Path
		targetQuery := c.Request.URL.RawQuery
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

package gateway

import (
	"github.com/gin-gonic/gin"
	"sort"
)

type HandlerOrderFunc struct {
	Order      int64
	FilterFunc gin.HandlerFunc
}

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerOrderFunc

// Global Filters
var globalFilters HandlersChain

func contextPathStripPrefixGlobalFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Request.URL.Path[len(contextPath):]
	}
}

func (globalFilter HandlerOrderFunc) Use() {
	globalFilters = append(globalFilters, globalFilter)
}

func loadGlobalFilters()  {
	HandlerOrderFunc{
		Order:      -100,
		FilterFunc: contextPathStripPrefixGlobalFilter(),
	}.Use()
	HandlerOrderFunc{
		Order:      -100,
		FilterFunc: contextPathStripPrefixGlobalFilter(),
	}.Use()
	sort.Sort(sort.Reverse(globalFilters))
}

// 重写 Len() 方法
func (a HandlersChain) Len() int {
	return len(a)
}

// 重写 Swap() 方法
func (a HandlersChain) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// 重写 Less() 方法， 从大到小排序
func (a HandlersChain) Less(i, j int) bool {
	return a[j].Order < a[i].Order
}


func initGlobalFilters(router *gin.Engine){
	loadGlobalFilters()
	for _,v:=range globalFilters{
		router.Use(v.FilterFunc)
	}
}
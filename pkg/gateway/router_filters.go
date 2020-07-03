package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"sort"
)

type RouterHandlerFunc struct {
	Order      int64
	FilterFunc gin.HandlerFunc
}

// HandlersChain defines a HandlerFunc array.
type RouterHandlersChain []RouterHandlerFunc

// Global Filters
var routerFilters RouterHandlersChain

func (h RouterHandlerFunc) Use() {
	routerFilters = append(routerFilters, h)
}

// 重写 Len() 方法
func (a RouterHandlersChain) Len() int {
	return len(a)
}

// 重写 Swap() 方法
func (a RouterHandlersChain) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// 重写 Less() 方法， 从大到小排序
func (a RouterHandlersChain) Less(i, j int) bool {
	return a[j].Order < a[i].Order
}

func init() {
	routerFilters = make(RouterHandlersChain, 0)
}

//初始化router对应的filterChain
func initRouterFiltersChain() {
	for _, r := range routerMap {
		tempFilterFunc := make([]RouterHandlerFunc, 0)
		for _, f := range r.Filters {
			order,filterFunc := getFilterFuncOrder(&f)
			if filterFunc != nil {
				tempFilterFunc = append(tempFilterFunc, RouterHandlerFunc{
					Order:            order,
					FilterFunc: filterFunc,
				})
			} else {
				fmt.Println("not find func", f.Name+"RouterFilterFunc")
			}
		}
		r.RouterFiltersChain = tempFilterFunc
		sort.Sort(sort.Reverse(r.RouterFiltersChain))
	}

}

//获取filter 执行顺序
func getFilterFuncOrder(f *Filter) (int64, gin.HandlerFunc) {
	methodName := f.Name + "RouterFilter"
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GetFilterFuncOrder faild", methodName)
		}
	}()
	//获取所有filter对象声明的以RouterFilter结尾的func
	methodValue := reflect.ValueOf(f).MethodByName(methodName)
	args := make([]reflect.Value, 1)
	args[0] = reflect.ValueOf(f.Name)
	results := methodValue.Call(args)
	order := results[0]
	rc := results[1]
	return order.Int(),rc.Interface().(gin.HandlerFunc)
}

//执行router filter Chain
//todo index
func executeRouterFiltersChain(router *Router, c *gin.Context) {
	//获取对应对router filtersChain
	index := 0
	for index < len(router.RouterFiltersChain) {
		router.RouterFiltersChain[index].FilterFunc(c)
		index++
	}
}

//测试RouterFilter ,返回值是filter 执行顺序，数值越小 越先执行
func (filter Filter) TestRouterFilter(key string) (int64,gin.HandlerFunc) {
	test := func(req *gin.Context)  {
	}
	return -1, test
}
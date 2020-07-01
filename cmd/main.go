package main

import (
	"fmt"
	"gin-gateway/pkg/gateway"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
)

func main() {
	var (
		port  = kingpin.Flag("port", "Server Port.").Default(viper.GetString("server.port")).String()
		debug = kingpin.Flag("debug", "Enabled Debug").Bool()
	)
	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}
	kingpin.Parse()
	quit := make(chan struct{})
	defer close(quit)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			select { //nolint: megacheck
			case <-sig:
				quit <- struct{}{}
			}
			os.Exit(0)
		}
	}()
	engine:=gin.Default()
	gateway.InitRouter(engine)
	info := fmt.Sprintf("%s:%s", "0.0.0.0", *port)
	engine.Run(info)

}

func init() {
	viper.SetConfigName("config")  //  设置配置文件名 (不带后缀)
	viper.AddConfigPath("configs") // 第一个搜索路径
	err := viper.ReadInConfig()    // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

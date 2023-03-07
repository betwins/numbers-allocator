package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maczh/mgin"
	"github.com/maczh/mgin/config"
	"github.com/maczh/mgin/i18n"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgrabbit"
	"net/http"
	_ "numbers-allocator/docs"
	"numbers-allocator/router"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

//@title	numbers-allocator
//@version 	1.0.0(numbers-allocator)
//@description	ecloud idrange

// 初始化命令行参数
func parseArgs() string {
	var configFile string
	flag.StringVar(&configFile, "f", os.Args[0]+".yml", "yml配置文件名")
	flag.Parse()
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if !strings.Contains(configFile, "/") {
		configFile = path + "/" + configFile
	}
	return configFile
}

func main() {
	//初始化配置，自动连接数据库和Nacos服务注册
	configFile := parseArgs()
	mgin.Init(configFile)
	i18n.Init()
	mgin.MGin.Use("rabbitmq", mgrabbit.Rabbit.Init, mgrabbit.Rabbit.Close, nil)

	//GIN的模式，生产环境可以设置成release
	gin.SetMode("debug")
	engine := router.SetupRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Config.App.Port),
		Handler: engine,
	}

	fmt.Println("|-----------------------------------|")
	fmt.Println("|   numbers-allocator  1.0.0    |")
	fmt.Println("|-----------------------------------|")
	fmt.Println("|  Go Http Server Start Successful  |")
	fmt.Println("|    Port:" + config.Config.GetConfigString("go.application.port") + "     Pid:" + fmt.Sprintf("%d", os.Getpid()) + "        |")
	fmt.Println("|-----------------------------------|")
	fmt.Println("")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error("HTTP server listen: " + err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalChan
	logs.Error("Get Signal:" + sig.String())
	logs.Error("Shutdown Server ...")
	mgin.MGin.SafeExit()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logs.Error("Server Shutdown:" + err.Error())
	}
	logs.Error("Server exiting")

}

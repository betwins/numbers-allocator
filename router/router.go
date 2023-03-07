package router

import (
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/maczh/mgin/i18n"
	"github.com/maczh/mgin/logs"
	"github.com/maczh/mgin/middleware/cors"
	"github.com/maczh/mgin/middleware/postlog"
	"github.com/maczh/mgin/middleware/trace"
	"github.com/maczh/mgin/middleware/xlang"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"numbers-allocator/errcode"
	"runtime/debug"
)

/*
*
统一路由映射入口
*/
func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	engine := gin.Default()

	//添加跟踪日志
	engine.Use(trace.TraceId())

	//设置接口日志
	engine.Use(postlog.RequestLogger())
	//添加跨域处理
	engine.Use(cors.Cors())

	//添加国际化处理
	engine.Use(xlang.RequestLanguage())

	//添加swagger支持
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//处理全局异常
	engine.Use(nice.Recovery(recoveryHandler))

	//设置404返回的内容
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, i18n.Error(404, "404"))
	})

	initAllocatorRouter(engine)

	return engine
}

func recoveryHandler(c *gin.Context, err any) {
	logs.Error("panic {} {}", err, string(debug.Stack()))
	c.JSON(http.StatusOK, errcode.UrlNotFound.MGError())
}

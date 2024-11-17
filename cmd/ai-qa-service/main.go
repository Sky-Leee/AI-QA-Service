package main

import (
	"ai-qa-service/internal/conf"
	"ai-qa-service/internal/coze"
	"ai-qa-service/internal/handler"
	"ai-qa-service/internal/job"
	"ai-qa-service/internal/logger"
	"ai-qa-service/internal/middleware"
	"ai-qa-service/internal/utils"
	"ai-qa-service/pkg/models"
	docs "ai-qa-service/pkg/docs"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	path := flag.String("c", "./config/config.json", "config path(file must be named 'config.json')")
	flag.Parse()

	conf.InitConf(*path)
	
	logger.InitLogger()

	utils.InitToken()

	models.InitMySQL()
	logger.Infof("Initializing MySQL successfully")

	coze.InitCoze()
	
	job.InitWorkers()
}

//	@title			AI-QA-Service 接口文档
//	@version		1.0
//	@description	包含了 AI-QA-Service 项目提供的接口

//	@contact.name	skylee
//	@contact.email	1350650238@qq.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		127.0.0.1:1145
// @BasePath	/api/v1
func main() {
	srv := GetSrv()

	idleConnsClosed := make(chan interface{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint // 阻塞，直到 SIGINT 信号产生

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt64("server.shutdown_waitting_time")))
		defer cancel()
		logger.Infof("Shutting down HTTP Server(wait for all connections to be closed)...")

		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf("Server shutdown failed: %v", err)
		}
		logger.Infof("Server closed successfully")
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("HTTP server ListenAndServe failed: %v", err)
	}

	<-idleConnsClosed // 直到 close 后，主线程才会退出
	logger.Infof("Waitting for all background tasks to complete...")
	job.Wait() // 等待所有后台任务结束才退出
	logger.Infof("Done.\n\nServer closed successfully")
}

func GetSrv() *http.Server {
	router := gin.New()

	// 中间件
	frontendPath := viper.GetString("router.corf.frontend_path")
	middlewares := []gin.HandlerFunc{logger.GinLogger(), logger.GinRecovery(true), middleware.CORF(frontendPath)}
	router.Use(middlewares...)

	// Swagger 接口文档
	if viper.GetBool("service.swagger.enable") {
		docs.SwaggerInfo.BasePath = "/api/v1"
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 注册路由
	v1 := router.Group("/api/v1")
	{
		usrGrp := v1.Group("/user")
		quesGrp := v1.Group("/question")
		statGrp := v1.Group("/statistics")
		
		quesGrp.Use(middleware.Auth(), middleware.VerifyToken())
		statGrp.Use(middleware.Auth(), middleware.VerifyToken())
		
		handler.RegistUserHandler(usrGrp)
		handler.RegistQuestionHandler(quesGrp)
		handler.RegistStatisticsHandler(statGrp)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("server.ip"), viper.GetInt("server.port")),
		Handler: router,
	}

	return srv
}

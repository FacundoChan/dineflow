package server

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunHTTPServer(serviceName string, wrapper func(r *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		// TODO: Warning log
		logrus.Panic("addr can't be empty")
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

func RunHTTPServerOnAddr(addr string, wrapper func(r *gin.Engine)) {
	apiRouter := gin.New()

	apiRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	wrapper(apiRouter)
	apiRouter.Group("api/")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

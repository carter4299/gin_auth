package myserverconfig

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var dev_logger = DevLogger()

func InitCors(r *gin.Engine) *gin.Engine {
	dev_logger.Info("Entered InitCors ()")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://%s:%s", os.Getenv("WEB_PROXY_HOST"), os.Getenv("WEB_PROXY_PORT"))},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Cookie", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	dev_logger.Info("Cors initialized ...")
	return r
}

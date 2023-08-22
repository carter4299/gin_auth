package router

import (
	auth "github.com/carter4299/gin_auth/auth"
	devlogger "github.com/carter4299/gin_auth/my_server_config"
	"github.com/gin-gonic/gin"
)

var dev_logger = devlogger.DevLogger()

func AuthRoutes(c *gin.Context) {
	dev_logger.Info("Entered AuthRoutes() in router / router.go ...")
	dev_logger.Info(c.FullPath())

	switch c.FullPath() {
	case "/api/go/auth/login":
		auth.Login(c)
	case "/api/go/auth/signup":
		auth.Signup(c)
	default:
		dev_logger.Error("Invalid path")
		c.JSON(400, gin.H{"error": "Invalid path"})
	}
}

func OtherRoutes(c *gin.Context) {
	dev_logger.Info("OtherRoutes")

	switch c.FullPath() {
	case "/api/go/other/valid_token":
		dev_logger.Debug("Passed Middleware")
		c.JSON(200, gin.H{"status": "success"})
	case "/api/go/other/logout":
		dev_logger.Debug("Passed Middleware")
		c.JSON(200, gin.H{"status": "success"})
	default:
		dev_logger.Error("Invalid path")
		c.JSON(400, gin.H{"error": "Invalid path"})
	}
}

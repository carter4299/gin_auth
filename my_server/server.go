package myserver

import (
	"fmt"
	"os"

	middleware "github.com/carter4299/gin_auth/middleware"
	myserverconfig "github.com/carter4299/gin_auth/my_server_config"
	router "github.com/carter4299/gin_auth/router"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)


func MyServer() {
    fmt.Println("Hello ! Good luck with your project ...")
    gin.SetMode(gin.DebugMode)

    r := gin.New()
    r = myserverconfig.InitCors(r)
    r.Use(myserverconfig.Logger(), gin.Recovery())
    r = middleware.InitMiddleware(r)

    r.Static("/assets", "./my_server/assets") 
    r.StaticFile("/", "./my_server/index.html") 
    r.NoRoute(func(c *gin.Context) {
        c.File("./my_server/index.html")
    })


    auth_routes := r.Group("/api/go/auth")
    {
        auth_routes.POST("/login", router.AuthRoutes)
        auth_routes.POST("/signup", router.AuthRoutes)
    }
    other_routes := r.Group("/api/go/other")
    {
        other_routes.POST("/logout", router.AuthRoutes)
        other_routes.GET("/valid_token", router.OtherRoutes)
    }


    r.Run(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
}

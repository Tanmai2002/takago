package server

import (
	"github.com/gin-gonic/gin"
)

func InitializeServer() *gin.Engine {
	r := gin.Default()
	AddDeploymentRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	return r
}

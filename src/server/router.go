package server

import (
	"matchamking/server/command"

	"github.com/gin-gonic/gin"
)

func getRouter(commands command.Commander) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	c := &command.UserAdd{}
	v1.POST("/"+c.Name(), commands.Register(c))
	return router
}

package server

import (
	"matchamking/src/server/command"

	"github.com/gin-gonic/gin"
)

func getRouter(commands command.Commander) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	var c command.ICommand
	c = &command.UserAdd{}
	v1.POST("/"+c.Name(), commands.Register(c))
	c = &command.UsersList{}
	v1.POST("/"+c.Name(), commands.Register(c))
	return router
}

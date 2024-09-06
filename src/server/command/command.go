package command

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ILogger interface {
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type ICommand interface {
	Name() string
	Apply() (any, error)
	Parse(ctx *gin.Context) error
}

type Commander interface {
	Register(command ICommand) gin.HandlerFunc
}

func NewManager(logger ILogger) *Manager {
	return &Manager{logger: logger}
}

type Manager struct {
	logger ILogger

	commands map[string]ICommand
}

func (c *Manager) Register(command ICommand) gin.HandlerFunc {
	if c.commands == nil {
		c.commands = make(map[string]ICommand)
	}
	c.commands[command.Name()] = command
	return func(ctx *gin.Context) {
		c.logger.Info("new request to command: %s", command.Name())
		if err := command.Parse(ctx); err != nil {
			c.logger.Error("cannot parse request: %s, error: %s", command.Name(), err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, err := command.Apply()
		if err != nil {
			c.logger.Error("cannot apply request: %s, error: %s", command.Name(), err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}

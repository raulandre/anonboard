package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raulandre/anonboard/config"
)

type Router interface {
	gin.IRouter
	Serve() error
}

type router struct {
	*gin.Engine
	c *config.Config
}

func NewRouter(c *config.Config) Router {
	config := c.Get()

	if config.GetString("ENVIROMENT") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if config.GetBool("app.log") {
		r.Use(gin.Logger())
	}

	return &router{
		c:      c,
		Engine: r,
	}
}

func (r *router) Serve() error {
	port := r.c.Get().GetString("app.port")
	return r.Run(":" + port)
}

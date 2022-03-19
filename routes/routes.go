package routes

import "github.com/raulandre/anonboard/controllers"

func (r *router) RegisterThreadRoutes(c controllers.ThreadController) {
	rg := r.Group("/api/threads")
	rg.GET("/", c.ListThreads)
	rg.POST("/", c.CreateThread)
	rg.GET("/:id", c.GetThread)
	rg.PUT("/:id", c.Report)
	rg.DELETE("/:id", c.Delete)
}

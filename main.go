package main

import (
	"github.com/raulandre/anonboard/config"
	"github.com/raulandre/anonboard/controllers"
	"github.com/raulandre/anonboard/database"
	"github.com/raulandre/anonboard/routes"
	"github.com/raulandre/anonboard/services"
)

func main() {
	c := config.NewConfig()
	r := routes.NewRouter(c)

	conn := database.NewDbConnection(c)

	ts := services.NewThreadService(conn)
	tc := controllers.NewThreadController(ts)

	r.RegisterThreadRoutes(tc)

	r.Serve()
}

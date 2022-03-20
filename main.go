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
	rs := services.NewReplyService(conn)

	tc := controllers.NewThreadController(ts)
	rc := controllers.NewReplyController(rs)

	r.RegisterThreadRoutes(tc)
	r.RegisterReplyRoutes(rc)

	r.Serve()
}

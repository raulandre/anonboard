package main

import (
	"fmt"

	"github.com/raulandre/anonboard/config"
	"github.com/raulandre/anonboard/database"
	"github.com/raulandre/anonboard/routes"
)

func main() {
	c := config.NewConfig()
	r := routes.NewRouter(c)
	conn := database.NewDbConnection(c)

	fmt.Println(conn)

	r.Serve()
}

package main

import (
	"github.com/raulandre/anonboard/config"
	"github.com/raulandre/anonboard/routes"
)

func main() {
	c := config.NewConfig()
	r := routes.NewRouter(c)

	r.Serve()
}

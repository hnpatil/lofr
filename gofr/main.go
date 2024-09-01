package main

import (
	"github.com/hnpatil/lofr"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.POST("/users", lofr.Handler(postUser))
	app.GET("/users/{id}", lofr.Handler(getUser))
	app.DELETE("/users/{id}", lofr.Handler(deleteUser))

	app.Run()
}

# lofr
[![Maintainability](https://api.codeclimate.com/v1/badges/4e76134aa21a7f79949e/maintainability)](https://codeclimate.com/github/hnpatil/lofr/maintainability)

lofr lets you write simpler [gofr](https://gofr.dev/) handlers.

**key features**
- Binds request body, query params and path variable into single request object.
- Reduces boilerplate code.

## Getting started
Add lofr to your project
````bash
go get github.com/hnpatil/lofr
````
Register paths with lofr Handler 
````go
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
````
Define handler input objects 
````go
package main

type PostUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age" default:"21"`
}

type GetUserRequest struct {
	ID    int    `path:"id"`
	Name  string `query:"name"`
	Email string `query:"email"`
	Age   int    `query:"age" default:"21"`
}

type DeleteUserRequest struct {
	ID int `path:"id"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
````
Implement lofr compatible handlers
````go
package main

import (
	"gofr.dev/pkg/gofr"
)

func postUser(ctx *gofr.Context, user *PostUserRequest) (*User, error) {
	// handle create user logic  
}

func getUser(ctx *gofr.Context, req *GetUserRequest) (*User, error) {
	// handle get user logic
}

func deleteUser(ctx *gofr.Context, req *DeleteUserRequest) error {
	// handle delete user logic
}
````

package main

import (
	"gofr.dev/pkg/gofr"
	"math/rand"
)

type UserID struct {
	ID int `path:"id"`
}

type PostUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age" default:"21"`
}

func postUser(ctx *gofr.Context, user *PostUserRequest) (*User, error) {
	return &User{
		ID:    rand.Int(),
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

type GetUserRequest struct {
	UserID
	Name  string `query:"name"`
	Email string `query:"email"`
	Age   int    `query:"age" default:"21"`
}

func getUser(ctx *gofr.Context, req *GetUserRequest) (*User, error) {
	name := "user name"
	email := "user@email.com"

	if req.Name != "" {
		name = req.Name
	}

	if req.Email != "" {
		email = req.Email
	}

	return &User{
		ID:    req.ID,
		Name:  name,
		Email: email,
		Age:   req.Age,
	}, nil
}

type DeleteUserRequest struct {
	UserID
}

func deleteUser(ctx *gofr.Context, req *DeleteUserRequest) error {
	ctx.Infof("deleted user with id %d", req.ID)
	return nil
}

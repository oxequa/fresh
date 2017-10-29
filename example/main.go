package main

import (
"github.com/tockins/fresh"
"net/http"
"github.com/tockins/fresh/example/models"
)



func main() {
	f := fresh.New()
	f.Config().SetPort(8080)
	// group
	g := f.Group("/todos/").Before(filter).After(filter)
	g.GET("/", list)
	g.GET("/{todoUuid}", single)
	g.GET("/{todoUuid}/users/{userUuid}", single)
	f.Run()
}

func list(f fresh.Context) error {
	data := []models.Todo{{Title: "Buy milk"}, {Title: "Car wash"}}
	return f.Response().JSON(http.StatusOK, data)
}

func single(f fresh.Context) error {
	data := models.Todo{
		Uuid: f.Request().URLParam("todoUuid"),
		Title: "Buy milk",
		UserUuid:  f.Request().URLParam("userUuid"),
		}
	return f.Response().JSON(http.StatusOK, data)
}

func filter(f fresh.Context) error{
	return nil
}


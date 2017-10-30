package main

import (
	"github.com/tockins/fresh"
	"net/http"
	"github.com/tockins/fresh/example/models"
	"golang.org/x/net/websocket"
	"fmt"
)

func main() {
	f := fresh.New()
	f.Config().SetPort(8080)
	// group
	g := f.Group("/todos/").Before(filter).After(filter)
	g.GET("/", list)
	g.GET("/{todoUuid}", single)
	g.GET("/{todoUuid}/users/{userUuid}", single)
	f.WS("ws",socket)
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


func socket(f fresh.Context) error{
	for {
		ws := f.Request().WS()
		// Write
		err := websocket.Message.Send(ws, "Hello, Client!")
		if err != nil {
			fmt.Println("E",err)
		}
		//
		//// Read
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Println("E1",err)
		}
		fmt.Println("msg", msg)
	}
	return nil
}

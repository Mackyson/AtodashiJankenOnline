package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

type Player struct {
	OwnRoom *Room
	Ws      *websocket.Conn
	Name    string
}

func NewPlayer(ownRoom *room, ws *websocket.Conn, name string) {
	return &Player{
		OwnRoom: ownRoom,
		Ws:      ws,
		Name:    name,
	}
}

type Room struct {
	Players *[2]Player
	Id      string
}

func NewRoom(id string) {
	return &Room{
		Players: [2]Player{},
		Id:      id,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/:greet", func(c *gin.Context) {
		hello := c.Param("greet")
		c.JSON(http.StatusOK, gin.H{
			"greet": hello,
		})
	})
	router.Run(":" + port)
}

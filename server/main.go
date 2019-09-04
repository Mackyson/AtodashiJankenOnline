package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

type Hand int

const (
	Gu Hand = iota
	Chi
	Pa
)

type Message struct {
	Type    string
	Side    string
	payload string
}

type Player struct {
	OwnRoom *Room
	Ws      *websocket.Conn
	Side    string
	Name    string
}

func (p *Player) listen() {
}

func NewPlayer(ownRoom *Room, ws *websocket.Conn, side string, name string) *Player {
	return &Player{
		OwnRoom: ownRoom,
		Ws:      ws,
		Side:    side,
		Name:    name,
	}
}

type Room struct {
	Players      [2]*Player
	IsReadyCh    chan *Message
	ChangeHandCh chan *Message
	ResultCh     chan *Message
	Id           string
}

func (r *Room) run() {
	for _, i := range r.Players {
		i.Ws.WriteJSON("{\"text\":\"Yee\"}")
	}

}

func NewRoom(id string) *Room {
	return &Room{
		Players:      [2]*Player{},
		IsReadyCh:    make(chan *Message, 1),
		ChangeHandCh: make(chan *Message, 1),
		ResultCh:     make(chan *Message, 1),
		Id:           id,
	}
}

func matching(playerCh chan *websocket.Conn) {
	for {
		p1Ch := <-playerCh
		p2Ch := <-playerCh
		r := NewRoom("a")
		p1 := NewPlayer(r, p1Ch, "Left", "1")
		p2 := NewPlayer(r, p2Ch, "Right", "2")
		r.Players[0] = p1
		r.Players[1] = p2
		go p1.listen()
		go p2.listen()
		go r.run()
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
	router.Static("/assets", "../static/")
	router.GET("/greet/:greet", func(c *gin.Context) {
		hello := c.Param("greet")
		c.JSON(http.StatusOK, gin.H{
			"greet": hello,
		})
	})
	var playerCh = make(chan *websocket.Conn, 2)
	router.GET("/match", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("%+v", err)
		}
		playerCh <- conn
	})
	go matching(playerCh)
	router.Run(":" + port)
}

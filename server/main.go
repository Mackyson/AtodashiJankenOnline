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
	Type    string `json:"type"`
	Side    string `json:"side"`
	payload string `json:"payload"`
}

type Player struct {
	OwnRoom *Room
	Ws      *websocket.Conn
	Side    string
	Name    string
}

func (p *Player) listen() {
	for {
		var msg *Message
		err := p.Ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("+v", err)
			//			p.OwnRoom.Disconnected(p)
			_ = p.Ws.Close()
			return
		} else {
			log.Printf("+v", msg.Type)
			switch msg.Type {
			case "isReady":
				p.OwnRoom.IsReadyCh <- msg
			case "changeHand":
				p.OwnRoom.ChangeHandCh <- msg
			case "result":
				p.OwnRoom.ResultCh <- msg
			}
		}
	}
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
	//var startCh = make(chan bool)
	for {
		select {
		//case msg := <-r.IsReadyCh:
		case <-r.IsReadyCh: //動作チェック
			r.do1()
		//case hand := <-r.ChangeHandCh:
		case <-r.ChangeHandCh:
			r.do2()
		//case result := <-r.ResultCh:
		case <-r.ResultCh:
			r.do3()
		}
	}
}
func (r *Room) do1() { //provisional function
	for _, p := range r.Players {
		p.Ws.WriteJSON("\"do\":\"1\"")
	}
}
func (r *Room) do2() { //provisional function
	for _, p := range r.Players {
		p.Ws.WriteJSON("\"do\":\"2\"")
	}
}
func (r *Room) do3() { //provisional function
	for _, p := range r.Players {
		p.Ws.WriteJSON("\"do\":\"3\"")
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

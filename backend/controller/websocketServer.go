package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	appsv1 "k8s.io/api/apps/v1"

)

type client struct {
	hub *Hub
	ws   *websocket.Conn
	send chan []byte
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var updateChannel = make(chan *appsv1.Deployment)

func StreamUpdateds(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	c := &client{
		hub: hub,
		send: make(chan []byte),
		ws:   ws,
	}

	c.hub.register <- c

	go c.writePump()
}

func (c *client) writePump() {
	defer func() {
		fmt.Printf("closing socket\n")
		c.ws.Close()
	}()

	for {
		select {
		case message := <-c.send:
			fmt.Printf("got msg: %s\n", string(message))
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Printf("c.ws.WriteMessage error \n")
				log.Println(err)
				return
			}
		}
	}
}


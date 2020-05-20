package controller

import "fmt"

type Hub struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client

	content string
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			fmt.Printf("registerd new client\n")
			h.clients[c] = true
		case c := <-h.unregister:
			fmt.Printf("unregisterd a client\n")
			if _, ok := h.clients[c]; ok {
				fmt.Printf("the hub is deleteting client\n")
				delete(h.clients, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- m:
					fmt.Printf("sending to the client channel\n")
					break
				default:
					fmt.Printf("the hub can't reach the client deleteting client\n")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
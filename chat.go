package main

import "fmt"

type Chat struct {
	rooms      map[string]*Room
	createRoom chan string
}

func NewChat() *Chat {
	return &Chat{
		rooms:      make(map[string]*Room),
		createRoom: make(chan string),
	}
}

func (c *Chat) run() {
	for {
		select {
		case room := <-c.createRoom:
			c.rooms[room] = NewRoom()
			go c.rooms[room].run()
		}
	}
}

type Room struct {
	clients   map[*Client]bool
	join      chan *Client
	leave     chan *Client
	broadcast chan *SendMessage
}

func NewRoom() *Room {
	return &Room{
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan *SendMessage),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			fmt.Println(r.clients)
		case client := <-r.leave:
			r.clients[client] = false
			fmt.Println(r.clients)
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

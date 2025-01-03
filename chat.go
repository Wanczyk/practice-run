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
	clients map[*Client]bool
	join    chan *Client
	leave   chan *Client
}

func NewRoom() *Room {
	return &Room{
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
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
		}
	}
}

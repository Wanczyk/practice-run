package src

import (
	"fmt"
)

type Chat struct {
	Rooms      map[string]*Room
	CreateRoom chan string
}

func NewChat() *Chat {
	return &Chat{
		Rooms:      make(map[string]*Room),
		CreateRoom: make(chan string),
	}
}

func (c *Chat) Run() {
	for {
		select {
		case room := <-c.CreateRoom:
			c.Rooms[room] = NewRoom()
			go c.Rooms[room].Run()
		}
	}
}

type Room struct {
	Clients   map[*Client]bool
	Join      chan *Client
	Leave     chan *Client
	Broadcast chan *SendMessage
}

func NewRoom() *Room {
	return &Room{
		Clients:   make(map[*Client]bool),
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Broadcast: make(chan *SendMessage),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
			fmt.Println(r.Clients)
		case client := <-r.Leave:
			if _, active := r.Clients[client]; active {
				delete(r.Clients, client)
				fmt.Println(r.Clients)
			}
		case message := <-r.Broadcast:
			for client, active := range r.Clients {
				if !active {
					continue
				}
				select {
				case client.Send <- message:
				default:
					delete(r.Clients, client)
				}
			}
		}
	}
}

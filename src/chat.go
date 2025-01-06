package src

import (
	"log"
	"sync"
)

type Chat struct {
	Rooms      map[string]*Room
	CreateRoom chan string
	mu         sync.RWMutex
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
			c.mu.Lock()
			if _, exists := c.Rooms[room]; !exists {
				log.Println("Creating room: ", room)
				c.Rooms[room] = NewRoom()
				go c.Rooms[room].Run()
			}
			c.mu.Unlock()
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
			log.Println("Client join: ", client)
			r.Clients[client] = true
		case client := <-r.Leave:
			if _, active := r.Clients[client]; active {
				log.Println("Client leave: ", client)
				delete(r.Clients, client)
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					log.Println("Cannot send message, deleting client: ", client)
					delete(r.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}

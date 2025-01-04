package src

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
			if _, exists := c.Rooms[room]; !exists {
				c.Rooms[room] = NewRoom()
				go c.Rooms[room].Run()
			}
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
		case client := <-r.Leave:
			if _, active := r.Clients[client]; active {
				delete(r.Clients, client)
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					delete(r.Clients, client)
					close(client.Send)
				}
			}
		}
	}
}

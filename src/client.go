package src

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Chat *Chat
	Conn *websocket.Conn
	Send chan *SendMessage
}

func ServeWs(chat *Chat, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{Chat: chat, Conn: conn, Send: make(chan *SendMessage)}

	go client.ReadJSON()
	go client.WriteJSON()
}

func (c *Client) ReadJSON() {
	defer func() {
		log.Println("Cannot read from client, closing connection.")
		c.Conn.Close()
	}()

	for {
		message := &IncomingMessage{}
		if err := c.Conn.ReadJSON(message); err != nil {
			return
		}

		if message.Data == nil {
			if err := c.handleErrorCommand(message.Command, "Data field is missing", 400); err != nil {
				return
			}
			continue
		}

		if err := c.processCommand(message); err != nil {
			return
		}
	}
}

func (c *Client) WriteJSON() {
	defer func() {
		log.Println("Cannot write JSON to client, closing connection.")
		c.Conn.Close()
	}()
	for {
		select {
		case message := <-c.Send:
			if message == nil {
				log.Println("Empty message")
				continue
			}
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Println("Error writing JSON:", err)
				return
			}
		}
	}
}

func (c *Client) processCommand(message *IncomingMessage) error {
	log.Printf("Incoming message: %+v\n", message)
	switch message.Command {
	case "create":
		c.Chat.CreateRoom <- message.Data.Room
		return c.handleSuccessCommand(message.Data.Room, CreateCommand, fmt.Sprintf("created room with name %s", message.Data.Room))
	case "join":
		return c.handleJoinCommand(message)
	case "leave":
		return c.handleLeaveCommand(message)
	case "send":
		return c.handleSendCommand(message)
	default:
		return c.handleErrorCommand(UnknownCommand, fmt.Sprintf("Command: %s not found", message.Command), 400)
	}
}

func (c *Client) handleJoinCommand(message *IncomingMessage) error {
	if room, ok := c.Chat.Rooms[message.Data.Room]; ok {
		room.Join <- c
		return c.handleSuccessCommand(message.Data.Room, JoinCommand, fmt.Sprintf("joined room with name %s", message.Data.Room))
	}
	return c.handleRoomNotFound(message.Data.Room, JoinCommand)
}

func (c *Client) handleLeaveCommand(message *IncomingMessage) error {
	if room, ok := c.Chat.Rooms[message.Data.Room]; ok && room.Clients[c] {
		room.Leave <- c
		return c.handleSuccessCommand(message.Data.Room, LeaveCommand, fmt.Sprintf("left room with name %s", message.Data.Room))
	}
	return c.handleRoomNotFound(message.Data.Room, LeaveCommand)
}

func (c *Client) handleSendCommand(message *IncomingMessage) error {
	if room, ok := c.Chat.Rooms[message.Data.Room]; ok && room.Clients[c] {
		mess := &SendMessage{Status: StatusSuccess, Command: SendCommand, Data: &Data{Room: message.Data.Room, Message: message.Data.Message}}
		room.Broadcast <- mess
		return nil
	}
	return c.handleRoomNotFound(message.Data.Room, SendCommand)
}

func (c *Client) handleRoomNotFound(room string, command Commands) error {
	return c.handleErrorCommand(command, fmt.Sprintf("room: %s not found", room), 404)
}

func (c *Client) handleSuccessCommand(room string, command Commands, message string) error {
	resMessage := &SendMessage{Status: StatusSuccess, Command: command, Data: &Data{Room: room, Message: message}}
	if err := c.Conn.WriteJSON(resMessage); err != nil {
		log.Println("Error writing JSON:", err)
		return err
	}
	return nil
}

func (c *Client) handleErrorCommand(command Commands, message string, code int) error {
	resMessage := &SendMessage{Status: StatusError, Command: command, Error: &Error{Code: code, Message: message}}
	if err := c.Conn.WriteJSON(resMessage); err != nil {
		log.Println("Error writing JSON:", err)
		return err
	}
	return nil
}

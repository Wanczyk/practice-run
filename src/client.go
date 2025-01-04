package src

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Chat *Chat
	Conn *websocket.Conn
	Send chan *SendMessage
}

func (c *Client) ReadJSON() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		message := &IncomeMessage{}
		err := c.Conn.ReadJSON(message)
		if err != nil {
			return
		}

		switch message.Command {
		case "create":
			c.Chat.CreateRoom <- message.Room
			resMessage := fmt.Sprintf("created room with name %s", message.Room)
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
				log.Println(err)
				return
			}
		case "join":
			if room, ok := c.Chat.Rooms[message.Room]; ok {
				room.Join <- c
				resMessage := fmt.Sprintf("joined room with name %s", message.Room)
				if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
					log.Println(err)
					return
				}
			} else {
				c.handleRoomNotFound(message.Room)
			}
		case "leave":
			if room, ok := c.Chat.Rooms[message.Room]; ok {
				room.Leave <- c
				resMessage := fmt.Sprintf("left room with name %s", message.Room)
				if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
					log.Println(err)
					return
				}
			} else {
				c.handleRoomNotFound(message.Room)
			}
		case "send":
			if room, ok := c.Chat.Rooms[message.Room]; ok {
				mess := &SendMessage{Room: message.Room, Message: message.Message}
				room.Broadcast <- mess
			} else {
				c.handleRoomNotFound(message.Room)
			}
		default:
			log.Println("unknown command")
		}
	}

}

func (c *Client) WriteJSON() {
	defer func() {
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message := <-c.Send:
			err := c.Conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) handleRoomNotFound(room string) {
	resMessage := fmt.Sprintf("room: %s not found", room)
	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
		log.Println(err)
	}
}

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	chat *Chat
	conn *websocket.Conn
	send chan *SendMessage
}

func (c *Client) ReadJSON() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		message := &IncomeMessage{}
		err := c.conn.ReadJSON(message)
		if err != nil {
			return
		}

		switch message.Command {
		case "create":
			c.chat.createRoom <- message.Room
			resMessage := fmt.Sprintf("created room with name %s", message.Room)
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
				log.Println(err)
				return
			}
		case "join":
			if room, ok := c.chat.rooms[message.Room]; ok {
				room.join <- c
				resMessage := fmt.Sprintf("joined room with name %s", message.Room)
				if err := c.conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
					log.Println(err)
					return
				}
			} else {
				log.Println("room not found")
			}
		case "leave":
			if room, ok := c.chat.rooms[message.Room]; ok {
				room.leave <- c
				resMessage := fmt.Sprintf("left room with name %s", message.Room)
				if err := c.conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
					log.Println(err)
					return
				}
			} else {
				log.Println("room not found")
			}
		case "send":
			if room, ok := c.chat.rooms[message.Room]; ok {
				mess := &SendMessage{Room: message.Room, Message: message.Message}
				room.broadcast <- mess
			} else {
				log.Println("room not found")
			}
		default:
			log.Println("unknown command")
		}
	}

}

func (c *Client) WriteJSON() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}
}

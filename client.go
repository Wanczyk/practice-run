package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	chat *Chat
	conn *websocket.Conn
}

func (c *Client) ReadJSON() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		message := &Message{}
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
				room.leave <- c
				resMessage := fmt.Sprintf("left room with name %s", message.Room)
				if err := c.conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
					log.Println(err)
					return
				}
			} else {
				log.Println("room not found")
			}
		default:
			log.Println("unknown command")
		}
	}

}

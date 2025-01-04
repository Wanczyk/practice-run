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
		c.Conn.Close()
	}()

	for {
		message := &IncomeMessage{}
		err := c.Conn.ReadJSON(message)
		if err != nil {
			return
		}

		if message.Data == nil {
			if err := c.handleErrorCommand(message.Command, "Data field is missing", 400); err != nil {
				return
			}
			continue
		}

		switch message.Command {
		case "create":
			c.Chat.CreateRoom <- message.Data.Room
			err = c.handleSuccessCommand(message.Data.Room, CreateCommand, fmt.Sprintf("created room with name %s", message.Data.Room))
		case "join":
			room, ok := c.Chat.Rooms[message.Data.Room]
			if ok {
				room.Join <- c
				err = c.handleSuccessCommand(message.Data.Room, JoinCommand, fmt.Sprintf("joined room with name %s", message.Data.Room))
			} else {
				err = c.handleRoomNotFound(message.Data.Room, JoinCommand)
			}
		case "leave":
			room, ok := c.Chat.Rooms[message.Data.Room]
			if ok {
				room.Leave <- c
				err = c.handleSuccessCommand(message.Data.Room, LeaveCommand, fmt.Sprintf("left room with name %s", message.Data.Room))
			} else {
				err = c.handleRoomNotFound(message.Data.Room, LeaveCommand)
			}
		case "send":
			room, ok := c.Chat.Rooms[message.Data.Room]
			if ok {
				mess := &SendMessage{Status: StatusSuccess, Command: SendCommand, Data: &Data{Room: message.Data.Room, Message: message.Data.Message}}
				room.Broadcast <- mess
			} else {
				err = c.handleRoomNotFound(message.Data.Room, SendCommand)
			}
		default:
			err = c.handleErrorCommand(UnknownCommand, fmt.Sprintf("Command: %s not found", message.Command), 400)
		}
		if err != nil {
			return
		}
	}

}

func (c *Client) WriteJSON() {
	defer func() {
		c.Conn.Close()
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

func (c *Client) handleRoomNotFound(room string, command Commands) error {
	return c.handleErrorCommand(command, fmt.Sprintf("room: %s not found", room), 404)
}

func (c *Client) handleSuccessCommand(room string, command Commands, message string) error {
	resMessage := &SendMessage{Status: StatusSuccess, Command: command, Data: &Data{Room: room, Message: message}}
	if err := c.Conn.WriteJSON(resMessage); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) handleErrorCommand(command Commands, message string, code int) error {
	resMessage := &SendMessage{Status: StatusError, Command: command, Error: &Error{Code: code, Message: message}}
	if err := c.Conn.WriteJSON(resMessage); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

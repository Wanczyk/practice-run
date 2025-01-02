package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var chat = &Chat{}

func main() {
	chat := NewChat()
	go chat.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chat, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(chat *Chat, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if messageType != websocket.TextMessage {
			log.Println("invalid message type")
			continue
		}
		command := strings.SplitN(string(raw), " ", 2)
		if len(command) < 2 {
			log.Println("invalid command format")
			continue
		}

		action, payload := command[0], command[1]

		switch action {
		case "create":
			chat.createRoom <- payload
			resMessage := fmt.Sprintf("created room with name %s", payload)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
				log.Println(err)
				return
			}
			if room, ok := chat.rooms[payload]; ok {
				room.join <- &Client{conn: conn}
				resMessage := fmt.Sprintf("joined room with name %s", payload)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
					log.Println(err)
					return
				}
			} else {
				log.Println("room not found")
			}
		case "join":
			if room, ok := chat.rooms[payload]; ok {
				room.join <- &Client{conn: conn}
				resMessage := fmt.Sprintf("joined room with name %s", payload)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(resMessage)); err != nil {
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

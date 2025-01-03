package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type IncomeMessage struct {
	Command string `json:"command"`
	Room    string `json:"room"`
	Message string `json:"message"`
}

type SendMessage struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}

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

	client := &Client{chat: chat, conn: conn, send: make(chan *SendMessage)}

	go client.ReadJSON()
	go client.WriteJSON()
}

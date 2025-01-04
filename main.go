package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"practice-run/src"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	chat := src.NewChat()
	go chat.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chat, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(chat *src.Chat, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &src.Client{Chat: chat, Conn: conn, Send: make(chan *src.SendMessage)}

	go client.ReadJSON()
	go client.WriteJSON()
}

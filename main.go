package main

import (
	"log"
	"net/http"
	"practice-run/src"
)

func main() {
	chat := src.NewChat()
	go chat.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		src.ServeWs(chat, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

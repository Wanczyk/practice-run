package tests

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"os"
	"practice-run/src"
	"reflect"
	"strings"
	"testing"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	chat := src.NewChat()
	go chat.Run()
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		src.ServeWs(chat, w, r)
	}))
	defer server.Close()

	code := m.Run()
	os.Exit(code)
}

func TestHandleWebSocket_HappyPath(t *testing.T) {
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not connect to WebSocket server: %v", err)
	}
	defer client.Close()

	message := &src.IncomingMessage{
		Command: "create",
		Data: &src.Data{
			Room: "test",
		},
	}
	if err = client.WriteJSON(message); err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	responseMessage := &src.SendMessage{}
	if err = client.ReadJSON(responseMessage); err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	expected := &src.SendMessage{Status: src.StatusSuccess, Command: src.CreateCommand, Data: &src.Data{Room: "test", Message: "created room with name test"}}

	if !reflect.DeepEqual(responseMessage, expected) {
		t.Errorf("expected response '%+v', but got '%+v'", message, responseMessage)
	}

	message = &src.IncomingMessage{
		Command: "join",
		Data: &src.Data{
			Room: "test",
		},
	}
	if err = client.WriteJSON(message); err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	responseMessage = &src.SendMessage{}
	if err = client.ReadJSON(responseMessage); err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	expected = &src.SendMessage{Status: src.StatusSuccess, Command: src.JoinCommand, Data: &src.Data{Room: "test", Message: "joined room with name test"}}

	if !reflect.DeepEqual(responseMessage, expected) {
		t.Errorf("expected response '%+v', but got '%+v'", message, responseMessage)
	}

	message = &src.IncomingMessage{
		Command: "send",
		Data: &src.Data{
			Room:    "test",
			Message: "hello",
		},
	}
	if err = client.WriteJSON(message); err != nil {
		t.Fatalf("could not write message to WebSocket server: %v", err)
	}

	responseMessage = &src.SendMessage{}
	if err = client.ReadJSON(responseMessage); err != nil {
		t.Fatalf("could not read message from WebSocket server: %v", err)
	}

	expected = &src.SendMessage{Status: src.StatusSuccess, Command: src.SendCommand, Data: &src.Data{Room: "test", Message: "hello"}}

	if !reflect.DeepEqual(responseMessage, expected) {
		t.Errorf("expected response '%+v', but got '%+v'", message, responseMessage)
	}
}

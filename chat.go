package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	store Store
}

func NewChatService(s *Store) *ChatService {
	return &ChatService{
		store: *s,
	}
}

func (s *ChatService) RegisterRouters(r *mux.Router) {
	r.HandleFunc("/chat/{id}", s.handleWebSocket)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var roomMap map[string]*room

func (s *ChatService) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomId := mux.Vars(r)["id"] // Gets params
	fmt.Println(roomId)
	var room *room
	var ok bool
	room, ok = roomMap[roomId]
	if !ok {
		room = newRoom()
		roomMap[roomId] = room
		go room.run()
	}
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket:  socket,
		receive: make(chan []byte, 256),
		room:    room,
	}
	room.join <- client
	defer func() { room.leave <- client }()
	go client.write()
	client.read()
}

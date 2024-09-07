package main

import "github.com/gorilla/websocket"

type client struct {
	socket *websocket.Conn
	receive chan []byte
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		// Mesajı room.forward kanalına gönderirken gönderen client'i de belirtiyoruz
		c.room.forward <- message{data: msg, sender: c}
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

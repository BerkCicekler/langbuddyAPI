package main

type room struct {
	// clients holds all current clients in this room.
	clients map[*client]bool

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// forward is a channel that holds incoming messages along with the sender.
	forward chan message
}

// message struct to hold the message data and the sender
type message struct {
	data   []byte
	sender *client
}

// newRoom creates a new chat room
func newRoom() *room {
	return &room{
		forward: make(chan message),  // forward channel now handles message struct
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true

		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)

		case msg := <-r.forward:
			// Send message to all clients except the sender
			for client := range r.clients {
				if client != msg.sender {
					select {
					case client.receive <- msg.data:
					default:
						delete(r.clients, client)
						close(client.receive)
					}
				}
			}
		}
	}
}
package src

import (
	"log"
	"rit/config"
)

// Hub maintains the set of assigned active clients and broadcasts messages to the clients.
type Hub struct {
	master       *Master
	number       uint64
	clients      map[*Client]bool
	broadcast    chan []byte
	register     chan *Client
	registerLock chan int
	unregister   chan *Client
}

func NewHub(master *Master) *Hub {
	return &Hub{
		master:       master,
		number:       master.hubCounter,
		broadcast:    make(chan []byte),
		register:     make(chan *Client, config.HubSize),
		registerLock: make(chan int, config.HubSize),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	log.Println("[HUB] Started")
	h.master.registerHubs <- h
	log.Println("[HUB] Sent registration of Hub to master")
	for {
		select {
		case client := <-h.register:
			h.registerLock <- 1
			client.setHub(h)
			h.clients[client] = true
			log.Printf("[HUB] Client assigned to Hub\n")
		case client := <-h.unregister:
			_ = <-h.registerLock
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

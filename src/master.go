package src

import (
	"fmt"
	"log"

	"rit/config"
	"strconv"
	"sync/atomic"
)

type Master struct {
	clientCounter   uint64
	hubCounter      uint64
	hubs            map[*Hub]uint64
	clients         map[uint64]*Client
	registerHubs    chan *Hub
	registerClients chan *Client
}

func NewMaster() *Master {
	return &Master{
		hubs:            make(map[*Hub]uint64),
		hubCounter:      0,
		clientCounter:   0,
		clients:         make(map[uint64]*Client),
		registerHubs:    make(chan *Hub),
		registerClients: make(chan *Client),
	}
}

func (m *Master) listenClientsRegistration() {
	log.Println("[MASTER] Start to listen Client registration")
	for {
		select {
		case client := <-m.registerClients:
			log.Printf("[HUB] Client %v registered\n", m.clientCounter)
			m.cookClient(client)
			freHubFound := m.assignClientToFreeHub(client)
			if !freHubFound {
				m.createNewHub(client)
			}
		}
	}
}

func (m *Master) cookClient(client *Client) {
	atomic.AddUint64(&m.clientCounter, 1)
	client.id = m.clientCounter
	m.clients[m.clientCounter] = client
}

func (m *Master) assignClientToFreeHub(client *Client) bool {
	freHubFound := false
	for hub := range m.hubs {
		if len(hub.registerLock) < config.HubSize {
			hub.register <- client
			freHubFound = true
		}
	}
	return freHubFound
}

func (m *Master) createNewHub(client *Client) {
	log.Println("[MASTER] All hubs are busy or not found, lets create new one")
	hub := NewHub(m)
	go hub.run()
	hub.register <- client
}

func (m *Master) listenHubRegistration() {
	log.Println("[MASTER] Start to listen Hub registration")
	for {
		select {
		case hub := <-m.registerHubs:
			log.Println("[MASTER] New hub came to be registered")
			m.cookHub(hub)
		}
	}
}

func (m *Master) cookHub(hub *Hub) {
	atomic.AddUint64(&m.hubCounter, 1)
	m.hubs[hub] = m.hubCounter
	log.Printf("[MASTER] Hub %v registered \n", m.hubCounter)
}

func (m *Master) Run() {
	go m.listenConsole()
	go m.listenHubRegistration()
	go m.listenClientsRegistration()
}

func (m *Master) listenConsole() {
	for {
		command, id, message, _ := parseCommand()
		log.Println("[MASTER] Master going to execute the command")
		switch command {
		case "send":
			m.sendBroadcast(message, id)
		case "sendc":
			m.sendPersonalMessage(message, id)
		}
	}
}

func (m *Master) sendPersonalMessage(message string, id uint64) {
	message = "personal: " + message
	//TODO this is the shitty temp fix
	id = id + 1
	_, ok := m.clients[id]
	if ok {
		log.Printf("[MASTER] Send Client %v message\n", id)
		m.clients[id].send <- []byte(message)
	}
}

func (m *Master) sendBroadcast(message string, id uint64) {
	message = "broadcast: " + message
	for hub := range m.hubs {
		if hub.number == id {
			log.Printf("[MASTER] Send Hub %v broadcast\n", id)
			hub.broadcast <- []byte(message)
		}
	}
}

// TODO: naive implementation
func parseCommand() (string, uint64, string, error) {
	var command, commandId, message string
	_, err := fmt.Scanln(&command, &commandId, &message)
	if err != nil {
		log.Println("Please, enter command & arguments")
	}
	id, _ := strconv.Atoi(commandId)
	return command, uint64(id) - 1, message, err
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	gws "github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./index.html") })
	http.HandleFunc("/ws-endpoint", ChatSocket)
	http.ListenAndServe(":8080", nil)
}

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Manager struct {
	clients map[*gws.Conn]bool
	lock    sync.Mutex
}

func (m *Manager) AddClient(c *gws.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.clients[c] = true
}

func (m *Manager) RemoveClient(c *gws.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.clients[c]; ok {
		delete(m.clients, c)
		c.Close()
	}
}

func (m *Manager) BroadcastMessage(message []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for c := range m.clients {
		if err := c.WriteMessage(gws.TextMessage, message); err != nil {
			log.Println(err)
			delete(m.clients, c)
			c.Close()
		}
	}
}

var manager = Manager{
	clients: make(map[*gws.Conn]bool),
}

func ChatSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	manager.AddClient(conn)
	defer manager.RemoveClient(conn)

	for {
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(clientMsg))

		rendered := []byte(fmt.Sprintf(`
                     	<div hx-swap-oob="beforeend:#response">
                       		<p>%s</p>
                	</div>

                `, clientMsg))

		manager.BroadcastMessage(rendered)

	}
}

package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type Hub struct {
	log *slog.Logger
	sync.Mutex
	conns      map[*Connection]bool
	register   chan *Connection
	unregister chan *Connection
	broadcast  chan []byte
}

type Connection struct {
	send   chan []byte
	ws     *websocket.Conn
	userId string
}

type Message struct {
	From string `json:"From"`
	To   string `json:"To"`
	Data string `json:"Data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHub(log *slog.Logger) *Hub {
	return &Hub{
		log:        log,
		conns:      make(map[*Connection]bool),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.Lock()
			h.conns[conn] = true
			h.Unlock()
			h.log.Info("Connection registered user: " + conn.userId)

		case conn := <-h.unregister:
			h.Lock()
			if _, ok := h.conns[conn]; ok {
				delete(h.conns, conn)
				close(conn.send)
				h.log.Info("Connection unregistered user: " + conn.userId)
			}
			h.Unlock()

		case message := <-h.broadcast:
			h.Lock()
			for conn := range h.conns {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(h.conns, conn)
				}
			}
			h.Unlock()
		}
	}
}

func (c *Connection) readLoop(h *Hub) {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			h.log.Error("Failed to read message " + err.Error())
			break
		}

		var msg Message
		if err = json.Unmarshal(message, &msg); err != nil {
			h.log.Error("Failed to unmarshal message " + err.Error())
			continue
		}

		h.log.Info("Message received from " + msg.From + " to " + msg.To + " data: " + msg.Data)

		for conn := range h.conns {
			if conn.userId == msg.To {
				select {
				case conn.send <- []byte(msg.Data):
				default:
					close(conn.send)
					delete(h.conns, conn)
				}
			}
		}
	}
}

func (c *Connection) writeLoop(h *Hub) {
	for msg := range c.send {
		h.log.Info("Writing message to " + c.userId + " data: " + string(msg))
		if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			h.log.Error("Failed to write message " + err.Error())
			break
		}
	}
}

func (h *Hub) HandlerWS(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Reg new connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("Failed to upgrade connection " + err.Error())
		return
	}

	id, err := r.Cookie("id")
	if err != nil {
		h.log.Error("Failed to get cookie " + err.Error())
		return
	}

	c := &Connection{
		send:   make(chan []byte, 256),
		ws:     conn,
		userId: id.Value,
	}

	h.register <- c
	go c.writeLoop(h)
	c.readLoop(h)
}

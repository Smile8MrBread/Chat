package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
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

//func (s *Server) HandleOrderBook(ws *websocket.Conn) {
//	s.log.Info("New incoming connection from order book: ", ws.RemoteAddr())
//
//	for {
//		d := fmt.Sprintf("Order book data => %s", time.Now().Format("2006-01-02T15:04:05 -07:00:00"))
//		ws.Write([]byte(d))
//		time.Sleep(time.Second * 2)
//	}
//}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.Lock()
			h.conns[conn] = true
			h.Unlock()
		case conn := <-h.unregister:
			if h.conns[conn] {
				h.Lock()
				delete(h.conns, conn)
				close(conn.send)
				h.Unlock()
			}
		case message := <-h.broadcast:
			for conn := range h.conns {
				conn.send <- message
			}
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
			if err == io.EOF {
				break
			}
			h.log.Error("Failed to read message", err)
			break
		}

		var msg Message
		if err = json.Unmarshal(message, &msg); err != nil {
			h.log.Error("Failed to read message", err)
			break
		}

		fmt.Println(msg)
		for conn := range h.conns {
			if conn.userId == msg.To {
				conn.send <- []byte(msg.Data)
				break
			}
		}

		//h.broadcast <- []byte(msg.Data)
	}
}

func (c *Connection) writeLoop(h *Hub) {
	for msg := range c.send {
		if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			h.log.Error("Failed to write message", err)
			break
		}
	}
}

//func (s *Server) HandleWs(ws *websocket.Conn, id int) {
//	s.log.Info("New incoming connection from client: ", ws.RemoteAddr())
//
//	s.Lock()
//	s.conns[id] = ws
//	s.Unlock()
//
//	s.readLoop(ws)
//}

func (h *Hub) HandlerWS(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Reg new connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("Failed to reg new connection", err)
		return
	}

	id, err := r.Cookie("id")
	if err != nil {
		h.log.Error("Failed to parse cookie", err)
		return
	}

	c := &Connection{
		send:   make(chan []byte),
		ws:     conn,
		userId: id.Value,
	}

	h.register <- c
	go c.writeLoop(h)
	c.readLoop(h)
}

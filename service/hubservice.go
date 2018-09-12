package service

import (
	"bytes"
	"fmt"
	"log"
	"resulguldibi/color-api/entity"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type SocketClient struct {
	hub *SocketHub

	// The websocket connection.
	conn *SocketConnection

	// Buffered channel of outbound messages.
	send chan []byte

	user entity.User
}

type SocketHub struct {

	// Registered clients.
	clients            map[*SocketClient]bool
	multiPlayMatches   map[*MultiPlayMatch]bool
	clientMatchMapping map[*SocketClient]*MultiPlayMatch

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *SocketClient

	// Unregister requests from clients.
	unregister chan *SocketClient

	unRegisterMultiPlay chan *SocketClient

	registerMultiPlay chan *SocketClient

	acceptMatchForMultiPlay chan *SocketClient

	registerMatch chan *MultiPlayMatch

	unRegisterMatch chan *MultiPlayMatch
}

type MultiPlayMatch struct {
	clients           map[*SocketClient]bool
	clientAcceptances map[*SocketClient]bool
	wg                sync.WaitGroup
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewSocketHub() *SocketHub {
	return &SocketHub{
		broadcast:               make(chan []byte),
		register:                make(chan *SocketClient),
		unregister:              make(chan *SocketClient),
		clients:                 make(map[*SocketClient]bool),
		clientMatchMapping:      make(map[*SocketClient]*MultiPlayMatch),
		multiPlayMatches:        make(map[*MultiPlayMatch]bool),
		registerMultiPlay:       make(chan *SocketClient),
		acceptMatchForMultiPlay: make(chan *SocketClient),
		unRegisterMultiPlay:     make(chan *SocketClient),
		registerMatch:           make(chan *MultiPlayMatch),
		unRegisterMatch:         make(chan *MultiPlayMatch),
	}
}

func (h *SocketHub) Broadcast() {

	for {
		select {

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

func (h *SocketHub) RegisterMatch() {

	for {
		select {

		case multiPlayMatch := <-h.registerMatch:

			h.multiPlayMatches[multiPlayMatch] = true
			fmt.Println("RegisterMatch")
			fmt.Println("len(multiPlayMatch.clients) ->", len(multiPlayMatch.clients))
			if multiPlayMatch != nil && multiPlayMatch.clients != nil && len(multiPlayMatch.clients) > 0 {
				for c, _ := range multiPlayMatch.clients {
					h.clientMatchMapping[c] = multiPlayMatch

					opponent := multiPlayMatch.GetOpponentOf(c)

					socketMessage := &SocketMessage{}
					socketMessage.Message = fmt.Sprintf("you matched with %s", opponent.user.Email)
					c.SendMessage(socketMessage)
				}
			}

			multiPlayMatch.wg.Done()
		}
	}
}

func (m *MultiPlayMatch) GetOpponentOf(client *SocketClient) *SocketClient {
	var opponent *SocketClient
	if m.clients != nil && len(m.clients) > 0 {
		for c, _ := range m.clients {
			if c != client {
				opponent = c
				break
			}
		}
	}
	return opponent
}

func (h *SocketHub) UnRegisterMatch() {

	for {
		select {

		case multiPlayMatch := <-h.unRegisterMatch:

			fmt.Println("UnRegisterMatch")
			fmt.Println("len(multiPlayMatch.clients) ->", len(multiPlayMatch.clients))
			if multiPlayMatch != nil && multiPlayMatch.clients != nil && len(multiPlayMatch.clients) > 0 {

				if _, ok := h.multiPlayMatches[multiPlayMatch]; ok {
					delete(h.multiPlayMatches, multiPlayMatch)
				}

				for c, _ := range multiPlayMatch.clients {

					if _, ok := h.clientMatchMapping[c]; ok {
						delete(h.clientMatchMapping, c)
					}

					socketMessage := &SocketMessage{}
					socketMessage.Message = "you unmatched !!!!"
					c.SendMessage(socketMessage)
				}
			}
		}
	}
}

func (h *SocketHub) UnRegister() {

	for {
		select {

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

func (h *SocketHub) Register() {

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		}
	}
}

func (h *SocketHub) AcceptMatchForMultiPlay() {

	for {
		select {

		case client := <-h.acceptMatchForMultiPlay:

			multiPlayMatch := h.clientMatchMapping[client]

			if multiPlayMatch != nil && multiPlayMatch.clients != nil && len(multiPlayMatch.clients) > 0 {

				if _, ok := multiPlayMatch.clients[client]; ok {

					multiPlayMatch.clientAcceptances[client] = true

					opponent := multiPlayMatch.GetOpponentOf(client)

					socketMessage := &SocketMessage{}
					socketMessage.Message = fmt.Sprintf("your opponent %s accepted matching with you ", client.user.Email)
					opponent.SendMessage(socketMessage)

					if multiPlayMatch.clientAcceptances != nil && len(multiPlayMatch.clientAcceptances) == 2 {

						for c, _ := range multiPlayMatch.clients {

							socketMessage := &SocketMessage{}
							socketMessage.Message = "you will start to playy !!!!"
							c.SendMessage(socketMessage)
						}
					}
				}
			}
		}
	}
}

func (h *SocketHub) RegisterMultiPlay() {

	clients := make(map[*SocketClient]bool)
	clientAcceptances := make(map[*SocketClient]bool)

	for {
		select {

		case client := <-h.registerMultiPlay:

			h.clients[client] = true

			if len(clients) < 2 {
				clients[client] = true
			}

			if len(clients) == 2 {

				newMatch := &MultiPlayMatch{}
				newMatch.clients = clients
				newMatch.clientAcceptances = clientAcceptances
				fmt.Println("before match ->")
				newMatch.wg.Add(1)
				h.registerMatch <- newMatch
				newMatch.wg.Wait()
				fmt.Println("after match ->")
				clients = make(map[*SocketClient]bool)
				clientAcceptances = make(map[*SocketClient]bool)
			}
		}
	}
}

func (h *SocketHub) UnRegisterMultiPlay() {

	for {
		select {

		case client := <-h.unRegisterMultiPlay:

			if match, ok := client.hub.clientMatchMapping[client]; ok {
				client.hub.unRegisterMatch <- match
			}
		}
	}
}

func serveWs(hub *SocketHub, ctx *gin.Context, user entity.User) {
	connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}

	socketConnection := &SocketConnection{connection: connection}

	client := &SocketClient{hub: hub, conn: socketConnection, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()

}

func (c *SocketClient) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.connection.Close()
	}()
	c.conn.connection.SetReadLimit(maxMessageSize)
	c.conn.connection.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.connection.SetPongHandler(func(string) error { c.conn.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *SocketClient) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

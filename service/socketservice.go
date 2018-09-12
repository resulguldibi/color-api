package service

import (
	"flag"
	"fmt"
	"resulguldibi/color-api/entity"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var socketClientPool = &SocketClientPool{items: make(map[string]*SocketClient)}

type SocketConnection struct {
	connection *websocket.Conn
}

type SocketMessage struct {
	Message string `json:"message"`
}

type SocketClientPool struct {
	sync.RWMutex
	items map[string]*SocketClient
}

func (pool *SocketClientPool) AddClient(id string, client *SocketClient) {
	pool.Lock()
	defer pool.Unlock()

	pool.items[id] = client
}

func (pool *SocketClientPool) GetClient(id string) *SocketClient {
	pool.Lock()
	defer pool.Unlock()

	return pool.items[id]
}

func (s *SocketService) CreateSocketConnection(ctx *gin.Context, user entity.User, hub *SocketHub) *SocketClient {

	connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}

	client := &SocketClient{hub: hub, conn: &SocketConnection{connection: connection}, send: make(chan []byte, 256), user: user}
	client.hub.register <- client

	socketClientPool.AddClient(user.Id, client)

	return client
}

func (s *SocketService) GetSocketConnection(ctx *gin.Context, user entity.User, hub *SocketHub) *SocketClient {

	return socketClientPool.GetClient(user.Id)

}

func (s *SocketService) AcceptMatchForMultiPlay(ctx *gin.Context, user entity.User, hub *SocketHub) {

	fmt.Println("Before GetSocketConnection")

	socketClient := s.GetSocketConnection(ctx, user, hub)

	fmt.Println("Before AcceptMatchForMultiPlay")

	socketClient.hub.acceptMatchForMultiPlay <- socketClient

	fmt.Println("After AcceptMatchForMultiPlay")
}

func (s *SocketService) RegisterForMultiPlay(ctx *gin.Context, user entity.User, hub *SocketHub) {

	socketClient := s.CreateSocketConnection(ctx, user, hub)

	fmt.Println("Before RegisterForMultiPlay")

	socketClient.hub.registerMultiPlay <- socketClient

	fmt.Println("After RegisterForMultiPlay")
}

func (s *SocketService) UnRegisterForMultiPlay(ctx *gin.Context, user entity.User, hub *SocketHub) {

	socketClient := s.GetSocketConnection(ctx, user, hub)

	fmt.Println("Before UnRegisterForMultiPlay")

	socketClient.hub.unRegisterMultiPlay <- socketClient

	fmt.Println("After UnRegisterForMultiPlay")
}

func (client *SocketClient) ReadMessage() {

	_, message, err := client.conn.connection.ReadMessage()
	if err != nil {
		panic(err)
	}

	socketMessage := &SocketMessage{}
	socketMessage.Message = fmt.Sprintf("you send me this message -> %s", string(message))

	client.SendMessage(socketMessage)
}

func (client *SocketClient) SendMessage(message *SocketMessage) {

	err := client.conn.connection.WriteJSON(message)

	if err != nil {
		panic(err)
	}
}

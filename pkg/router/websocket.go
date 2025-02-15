package router

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader  websocket.Upgrader
	handlers  map[string]WSHandlerFunc
	onConnect func(*WSConnection)
	onClose   func(*WSConnection)
}

type WSConnection struct {
	conn    *websocket.Conn
	handler *WebSocketHandler
	data    map[string]interface{}
	send    chan []byte
	mu      sync.RWMutex
}

type WSHandlerFunc func(*WSConnection, []byte)

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Configure as needed
			},
		},
		handlers: make(map[string]WSHandlerFunc),
	}
}

func (wsh *WebSocketHandler) HandleFunc(event string, handler WSHandlerFunc) {
	wsh.handlers[event] = handler
}

func (wsh *WebSocketHandler) OnConnect(fn func(*WSConnection)) {
	wsh.onConnect = fn
}

func (wsh *WebSocketHandler) OnClose(fn func(*WSConnection)) {
	wsh.onClose = fn
}

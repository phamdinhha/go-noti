package socket

import (
	"context"
	"fmt"
	"go-noti/models"
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"github.com/gin-gonic/gin"
)

type RealtimeNotiEngine interface {
	// GetUserSockets(userID int)
	EmitToRoom(room string, key string, data interface{}) error
	EmitToUser(userID string, key string, data interface{}) error
	// EmitToGroupOfUser(userGroup string, key string, data interface{}) error
	Run(ctx context.Context, router *gin.Engine) error
}

type rnEngine struct {
	server  *socketio.Server
	storage map[int][]AppSocket
	locker  *sync.RWMutex
}

func NewEngine() *rnEngine {
	return &rnEngine{
		storage: make(map[int][]AppSocket),
		locker:  new(sync.RWMutex),
	}
}

func (e *rnEngine) saveAppSocket(userID int, appSocket AppSocket) {
	e.locker.Lock()
	if v, ok := e.storage[userID]; ok {
		e.storage[userID] = append(v, appSocket)
	} else {
		e.storage[userID] = []AppSocket{appSocket}
	}
	e.locker.Unlock()
}

func (e *rnEngine) getAppSocket(userID int) []AppSocket {
	e.locker.RLock()
	defer e.locker.RUnlock()
	return e.storage[userID]
}

func (e *rnEngine) removeAppSocket(userID int, appSocket AppSocket) {
	e.locker.Lock()
	defer e.locker.Unlock()
	if v, ok := e.storage[userID]; ok {
		for i := range v {
			if v[i] == appSocket {
				e.storage[userID] = append(v[:i], v[i+1:]...)
				break
			}
		}
	}
}

func (e *rnEngine) UserSockets(userID int) []AppSocket {
	var sockets []AppSocket
	if socks, ok := e.storage[userID]; ok {
		return socks
	}
	return sockets
}

func (e *rnEngine) EmitToRoom(room string, key string, data interface{}) error {
	e.server.BroadcastToRoom("/", room, key, data)
	return nil
}

func (e *rnEngine) EmitToUser(userID int, key string, data interface{}) error {
	sockets := e.getAppSocket(userID)
	for _, s := range sockets {
		s.Emit(key, data)
	}
	return nil
}

func (e *rnEngine) Run(ctx context.Context, router *gin.Engine) error {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{websocket.Default},
	})

	e.server = server
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("On connect: ", s.ID(), " IP: ", s.RemoteAddr())
		s.SetContext("")
		return nil
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Meet error: ", e)
	})
	server.OnDisconnect("/", func(c socketio.Conn, s string) {
		fmt.Println("closed with reason: ", s)
	})
	server.OnEvent("/", "authenticate", func(s socketio.Conn, token string) {
		// verify token here
		// can do authorization here
		fmt.Println("On user authentication")
		user := models.NewFakeUser()
		appSock := NewAppSocket(s, user)
		e.saveAppSocket(user.ID, appSock)
		s.Emit("authenticated", user)
		// server.OnEvent("/", "UserUpdateLocation", )
	})
	go server.Serve()
	defer server.Close()
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("../asset"))
	return nil
}

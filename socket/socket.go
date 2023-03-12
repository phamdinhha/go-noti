package socket

import (
	"go-noti/pkg"
	"net"
	"net/http"
	"net/url"
)

type Conn interface {
	// ID returns session id
	ID() string
	Close() error
	URL() url.URL
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteHeader() http.Header

	// Context of this connection. You can save one context for one
	// connection, and share it between all handlers. The handlers
	// is called in one goroutine, so no need to lock context if it
	// only be accessed in one connection.
	Context() interface{}
	SetContext(v interface{})
	Namespace() string
	Emit(msg string, v ...interface{})

	// Broadcast server side apis
	Join(room string)
	Leave(room string)
	LeaveAll()
	Rooms() []string
}

type AppSocket interface {
	Conn
	pkg.Requester
}

type appSocket struct {
	Conn
	pkg.Requester
}

func NewAppSocket(conn Conn, requester pkg.Requester) *appSocket {
	return &appSocket{conn, requester}
}

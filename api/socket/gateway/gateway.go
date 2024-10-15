package gateway

import (
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/zishang520/socket.io/socket"
)

type socketGateway struct {
	server    *socket.Server
	namespace string
}

func NewSocketGateway(server *socket.Server, namespace string) domain.SocketGateway {
	return &socketGateway{
		server:    server,
		namespace: namespace,
	}
}

func (g *socketGateway) OnConnection(callback func(socketio *socket.Socket)) {
	g.server.Of(g.namespace, nil).On("connection", func(clients ...any) {
		socketio := clients[0].(*socket.Socket)
		callback(socketio)
	})
}

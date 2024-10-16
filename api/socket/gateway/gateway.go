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

func (g *socketGateway) EmitRoom(room, event string, data interface{}) {
	g.server.Of(g.namespace, nil).To(socket.Room(room)).Emit(event, data)
}

func (g *socketGateway) JoinRoom(socketio *socket.Socket, room string) {
	socketio.Join(socket.Room(room))
}

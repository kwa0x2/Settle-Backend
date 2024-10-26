package gateway

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/utils"
	"github.com/zishang520/socket.io/socket"
)

type socketGateway struct {
	server    *socket.Server
	namespace string
	env       *bootstrap.Env
}

func NewSocketGateway(server *socket.Server, namespace string, env *bootstrap.Env) domain.SocketGateway {
	return &socketGateway{
		server:    server,
		namespace: namespace,
		env:       env,
	}
}

func (g *socketGateway) OnConnection(callback func(socketio *socket.Socket)) {
	g.server.Of(g.namespace, nil).On("connection", func(clients ...any) {
		socketio := clients[0].(*socket.Socket)

		auth, ok := socketio.Handshake().Auth.(map[string]interface{})
		if !ok {
			fmt.Println("Authorization data is not valid")
			socketio.Disconnect(true)
			return
		}

		token, ok := auth["token"].(string)
		if !ok || token == "" {
			fmt.Println("Authorization token is missing")
			socketio.Disconnect(true)
			return
		}

		user, err := utils.IsAuthorized(token, g.env.AccessTokenSecret)
		if err != nil {
			fmt.Println("Authentication failed:", err)
			socketio.Disconnect(true)
			return
		}

		socketio.SetData(user)
		callback(socketio)
	})
}

func (g *socketGateway) EmitRoom(room, event string, data interface{}) {
	g.server.Of(g.namespace, nil).To(socket.Room(room)).Emit(event, data)
}

func (g *socketGateway) JoinRoom(socketio *socket.Socket, room string) {
	socketio.Join(socket.Room(room))
}

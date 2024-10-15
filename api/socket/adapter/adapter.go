package adapter

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/zishang520/socket.io/socket"
)

type socketAdapter struct {
	gateway domain.SocketGateway
}

func NewSocketAdapter(gateway domain.SocketGateway) domain.SocketAdapter {
	return &socketAdapter{
		gateway: gateway,
	}
}

func (adapter *socketAdapter) HandleConnection() {
	adapter.gateway.OnConnection(func(socketio *socket.Socket) {
		fmt.Println("connected")

	})
}

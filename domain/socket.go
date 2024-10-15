package domain

import "github.com/zishang520/socket.io/socket"

type SocketAdapter interface {
	HandleConnection()
}

type SocketGateway interface {
	OnConnection(callback func(socketio *socket.Socket))
}

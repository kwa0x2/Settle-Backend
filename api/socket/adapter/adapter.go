package adapter

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/zishang520/socket.io/socket"
)

type socketAdapter struct {
	gateway        domain.SocketGateway
	messageUsecase domain.MessageUsecase
}

func NewSocketAdapter(gateway domain.SocketGateway, messageUsecase domain.MessageUsecase) domain.SocketAdapter {
	return &socketAdapter{
		gateway:        gateway,
		messageUsecase: messageUsecase,
	}
}

func (adapter *socketAdapter) HandleConnection() {
	adapter.gateway.OnConnection(func(socketio *socket.Socket) {
		fmt.Println("connected")

		adapter.gateway.JoinRoom(socketio, "00000000-0000-0000-0000-000000000000")

		socketio.On("sendMessage", func(args ...any) {
			adapter.handleSendMessage(args...)
		})

		socketio.On("deleteMessage", func(args ...any) {
			adapter.handleDeleteMessage(args...)
		})

		socketio.On("editMessage", func(args ...any) {
			adapter.handleEditMessage(args...)
		})
	})

}

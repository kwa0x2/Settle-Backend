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
		user, ok := socketio.Data().(*domain.User)
		if !ok {
			fmt.Println("User not found in socket data")
			return
		}

		fmt.Println("Connected user:", user)

		adapter.gateway.JoinRoom(socketio, "000000000000000000000001")

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

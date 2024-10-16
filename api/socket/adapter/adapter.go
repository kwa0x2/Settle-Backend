package adapter

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
			fmt.Println("send message triggered")

			data, ok := args[0].(map[string]interface{})
			if !ok {
				utils.Log().Error(`socket message type error`)
				return
			}

			message := &domain.Message{
				Content:  data["content"].(string),
				SenderID: data["sender_id"].(string),
				RoomID:   data["room_id"].(string),
			}

			if attachmentData, ok := data["attachment"].(map[string]interface{}); ok {
				message.Attachment = &domain.Attachment{
					ID:          attachmentData["id"].(primitive.ObjectID),
					Filename:    attachmentData["filename"].(string),
					Size:        int(attachmentData["size"].(float64)),
					Url:         attachmentData["url"].(string),
					ContentType: attachmentData["content_type"].(string),
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
				}
			}

			if repliedMessageData, ok := data["replied_message"].(map[string]interface{}); ok {
				message.RepliedMessage = &domain.Message{
					ID:        repliedMessageData["id"].(primitive.ObjectID),
					Content:   repliedMessageData["content"].(string),
					SenderID:  repliedMessageData["sender_id"].(string),
					RoomID:    repliedMessageData["room_id"].(string),
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				}
			}

			err := adapter.messageUsecase.Create(message)
			if err != nil {
				fmt.Println("error creating message:", err)
				return
			}

			adapter.gateway.EmitRoom("00000000-0000-0000-0000-000000000000", "message", message)
		})

		socketio.On("deleteMessage", func(args ...any) {
			fmt.Println("delete message triggered")

			data, ok := args[0].(map[string]interface{})
			if !ok {
				utils.Log().Error(`socket message type error`)
				return
			}

			messageID, err := primitive.ObjectIDFromHex(data["id"].(string))
			if err != nil {
				fmt.Errorf("invalid message id: %v", err)
				return
			}

			err = adapter.messageUsecase.SoftDelete(messageID)
			if err != nil {
				fmt.Println("error creating message:", err)
				return
			}
		})
	})

}

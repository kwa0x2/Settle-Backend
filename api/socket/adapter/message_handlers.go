package adapter

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/zishang520/engine.io/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func (adapter *socketAdapter) handleSendMessage(args ...any) {
	data, ok := args[0].(map[string]interface{})
	if !ok {
		utils.Log().Error(`socket message type error`)
		return
	}

	roomID, err := uuid.Parse(data["room_id"].(string))
	if err != nil {
		utils.Log().Error(`invalid roomid format`)
		return
	}

	message := &domain.Message{
		Content:  data["content"].(string),
		SenderID: data["sender_id"].(string),
		RoomID:   roomID.String(),
	}

	if attachmentData, ok := data["attachment"].(map[string]interface{}); ok {
		message.Attachment = &domain.Attachment{
			ID:          attachmentData["id"].(bson.ObjectID),
			Filename:    attachmentData["filename"].(string),
			Size:        attachmentData["size"].(int64),
			Url:         attachmentData["url"].(string),
			ContentType: attachmentData["content_type"].(string),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
	}

	if repliedMessageData, ok := data["replied_message"].(map[string]interface{}); ok {
		message.RepliedMessage = &domain.Message{
			ID:        repliedMessageData["id"].(bson.ObjectID),
			Content:   repliedMessageData["content"].(string),
			SenderID:  repliedMessageData["sender_id"].(string),
			RoomID:    repliedMessageData["room_id"].(string),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
	}

	sendErr := adapter.SendMessage(message)
	if sendErr != nil {
		utils.Log().Error(sendErr.Error())
		return
	}

	utils.Log().Error("success", message)
}

func (adapter *socketAdapter) SendMessage(message *domain.Message) error {
	err := adapter.messageUsecase.Create(message)
	if err != nil {
		fmt.Println("error creating message:", err)
		return err
	}

	adapter.gateway.EmitRoom("00000000-0000-0000-0000-000000000000", "message", message)
	return nil
}

func (adapter *socketAdapter) handleDeleteMessage(args ...any) {
	data, ok := args[0].(map[string]interface{})
	if !ok {
		utils.Log().Error(`socket message type error`)
		return
	}

	messageID, err := bson.ObjectIDFromHex(data["id"].(string))
	if err != nil {
		fmt.Errorf("invalid message id: %v", err)
		return
	}

	err = adapter.DeleteMessage(messageID)
	if err != nil {
		utils.Log().Error(err.Error())
		return
	}

	utils.Log().Error("success", messageID)
}

func (adapter *socketAdapter) DeleteMessage(messageID bson.ObjectID) error {
	err := adapter.messageUsecase.SoftDelete(messageID)
	if err != nil {
		fmt.Println("error creating message:", err)
		return err
	}

	notifyData := map[string]interface{}{
		"message_id": messageID.Hex(),
	}

	adapter.gateway.EmitRoom("00000000-0000-0000-0000-000000000000", "delete_message", notifyData)
	return nil
}

func (adapter *socketAdapter) handleEditMessage(args ...any) {
	data, ok := args[0].(map[string]interface{})
	if !ok {
		utils.Log().Error(`socket message type error`)
		return
	}

	messageID, err := bson.ObjectIDFromHex(data["id"].(string))
	if err != nil {
		fmt.Errorf("invalid message id: %v", err)
		return
	}

	err = adapter.EditMessage(messageID, data["content"].(string))
	if err != nil {
		utils.Log().Error(err.Error())
		return
	}

	utils.Log().Error("success", messageID)
}

func (adapter *socketAdapter) EditMessage(messageID bson.ObjectID, content string) error {
	err := adapter.messageUsecase.EditMessage(messageID, content)
	if err != nil {
		fmt.Println("error editing message:", err)
		return err
	}

	notifyData := map[string]interface{}{
		"message_id":     messageID.Hex(),
		"edited_message": content,
	}

	adapter.gateway.EmitRoom("00000000-0000-0000-0000-000000000000", "edit_message", notifyData)
	return nil
}

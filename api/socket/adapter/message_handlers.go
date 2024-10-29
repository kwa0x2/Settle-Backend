package adapter

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/utils"
	su "github.com/zishang520/engine.io/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (adapter *socketAdapter) handleSendMessage(args ...any) {
	data, ok := args[0].(map[string]interface{})
	if !ok {
		su.Log().Error("Invalid message format: data should be a map")
		return
	}

	roomID, err := utils.ParseObjectIDFromData(data, "RoomID")
	if err != nil {
		su.Log().Error("Invalid room ID format")
		return
	}

	su.Log().Info("Parsing sender: %+v", data["Sender"])

	sender, err := utils.ParseUser(data["Sender"].(map[string]interface{}))
	if err != nil {
		su.Log().Error("Invalid sender format")
		return
	}

	message := &domain.Message{
		Content: utils.ExtractString(data, "Content"),
		Sender:  sender,
		RoomID:  roomID,
	}

	if attachmentData, ok := data["Attachment"].(map[string]interface{}); ok {
		attachment, err := utils.ParseAttachment(attachmentData)
		if err != nil {
			su.Log().Error(err.Error())
			return
		}
		message.Attachment = attachment
	}

	if repliedMessageData, ok := data["RepliedMessage"].(map[string]interface{}); ok {
		repliedMessage, err := utils.ParseRepliedMessage(repliedMessageData)
		if err != nil {
			su.Log().Error(err.Error())
			return
		}
		message.RepliedMessage = repliedMessage
	}

	if err := adapter.SendMessage(message); err != nil {
		su.Log().Error(err.Error())
		return
	}

	su.Log().Info("Message sent successfully", message)
}

func (adapter *socketAdapter) SendMessage(message *domain.Message) error {

	err := adapter.messageUsecase.CreateAndUpdateRoom(message)
	if err != nil {
		fmt.Println("error creating message:", err)
		return err
	}

	notifyData := map[string]interface{}{
		"Action":  "NewMessage",
		"Message": message,
	}

	adapter.gateway.EmitRoom(message.RoomID.Hex(), "Message", notifyData)
	return nil
}

func (adapter *socketAdapter) handleDeleteMessage(args ...any) {
	data, ok := args[0].(map[string]interface{})
	if !ok {
		su.Log().Error(`socket message type error`)
		return
	}

	roomID, err := utils.ParseObjectIDFromData(data, "RoomID")
	if err != nil {
		su.Log().Error("Invalid room ID format")
		return
	}

	messageID, err := utils.ParseObjectIDFromData(data, "MessageID")
	if err != nil {
		su.Log().Error(err.Error())
		return
	}

	err = adapter.DeleteMessage(messageID, roomID)
	if err != nil {
		su.Log().Error(err.Error())
		return
	}

	su.Log().Error("success", messageID)
}

func (adapter *socketAdapter) DeleteMessage(messageID, roomID bson.ObjectID) error {
	err := adapter.messageUsecase.SoftDelete(messageID)
	if err != nil {
		fmt.Println("error creating message:", err)
		return err
	}

	notifyData := map[string]interface{}{
		"Action":    "DeleteMessage",
		"MessageID": messageID.Hex(),
		"RoomID":    roomID.Hex(),
	}

	adapter.gateway.EmitRoom(roomID.Hex(), "Message", notifyData)
	return nil
}

func (adapter *socketAdapter) handleEditMessage(args ...any) {
	data, ok := args[0].(map[string]interface{})
	if !ok {
		su.Log().Error(`socket message type error`)
		return
	}

	roomID, err := utils.ParseObjectIDFromData(data, "RoomID")
	if err != nil {
		su.Log().Error("Invalid room ID format")
		return
	}

	messageID, err := utils.ParseObjectIDFromData(data, "MessageID")
	if err != nil {
		su.Log().Error(err.Error())
		return
	}

	content, ok := data["Content"].(string)
	if !ok {
		su.Log().Error("content must be a string")
		return
	}

	err = adapter.EditMessage(messageID, roomID, content)
	if err != nil {
		su.Log().Error(err.Error())
		return
	}

	su.Log().Error("success", messageID)
}

func (adapter *socketAdapter) EditMessage(messageID, roomID bson.ObjectID, content string) error {
	err := adapter.messageUsecase.EditMessage(messageID, content)
	if err != nil {
		fmt.Println("error editing message:", err)
		return err
	}

	notifyData := map[string]interface{}{
		"Action":     "EditMessage",
		"RoomID":     roomID.Hex(),
		"MessageID":  messageID.Hex(),
		"NewContent": content,
	}

	adapter.gateway.EmitRoom(roomID.Hex(), "Message", notifyData)
	return nil
}

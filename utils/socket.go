package utils

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"github.com/zishang520/engine.io/utils"
	socketUtils "github.com/zishang520/engine.io/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// region "ExtractArgs" extracts the data and callback function from socket arguments.
func ExtractArgs(args []any) (map[string]interface{}, func([]interface{}, error)) {
	// Check if there are enough arguments
	if len(args) < 2 {
		utils.Log().Error(`not enough arguments`) // Log an error if not enough arguments are provided
		return nil, nil
	}

	// Extract data from the first argument and check its type
	data, ok := args[0].(map[string]interface{})
	if !ok {
		utils.Log().Error(`socket message type error`) // Log an error if the data type is incorrect
		return nil, nil
	}

	// Extract the callback function from the second argument and check its type
	callback, ok := args[1].(func([]interface{}, error))
	if !ok {
		utils.Log().Error(`callback function type error`) // Log an error if the callback type is incorrect
		return nil, nil
	}

	return data, callback // Return the extracted data and callback function
}

// endregion

// region Response defines a structured response format for socket communication.
type Response struct {
	Status  string `json:"status"`  // Response status (success/error)
	Message string `json:"message"` // Response message
}

// endregion

// region "SendResponse" sends a structured response back through the callback
func SendResponse(callback func([]interface{}, error), status, message string) {
	response := []interface{}{Response{Status: status, Message: message}} // Create a response object
	callback(response, nil)                                               // Invoke the callback with the response
}

// endregion

// region "LogError" logs an error message and sends an error response
func LogError(callback func([]interface{}, error), message string) {
	socketUtils.Log().Error(message)         // Log the error message
	SendResponse(callback, "error", message) // Send an error response
}

// endregion

// region "LogSuccess" sends a success response
func LogSuccess(callback func([]interface{}, error), message string) {
	SendResponse(callback, "success", message) // Send a success response
}

// endregion

func ParseObjectIDFromData(data map[string]interface{}, key string) (bson.ObjectID, error) {
	idStr, ok := data[key].(string)
	if !ok {
		return bson.ObjectID{}, fmt.Errorf("%s must be a string", key)
	}

	objectID := bson.ObjectID{}
	objectID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return bson.ObjectID{}, fmt.Errorf("invalid ObjectID format for %s: %w", key, err)
	}

	return objectID, nil
}

func ExtractString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func ParseAttachment(data map[string]interface{}) (*domain.Attachment, error) {
	attachmentID, err := ParseObjectIDFromData(data, "id")
	if err != nil {
		return nil, fmt.Errorf("Invalid attachment ID format")
	}

	return &domain.Attachment{
		ID:          attachmentID,
		Filename:    ExtractString(data, "filename"),
		Size:        int64(data["size"].(float64)),
		Url:         ExtractString(data, "url"),
		ContentType: ExtractString(data, "content_type"),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func ParseRepliedMessage(data map[string]interface{}) (*domain.Message, error) {
	repliedMessageID, err := ParseObjectIDFromData(data, "id")
	if err != nil {
		return nil, fmt.Errorf("Invalid replied message ID format")
	}

	roomID, err := ParseObjectIDFromData(data, "room_id")
	if err != nil {
		return nil, fmt.Errorf("Invalid replied room ID format")
	}

	readStatus, err := parseReadStatus(data)
	if err != nil {
		return nil, err
	}

	return &domain.Message{
		ID:         repliedMessageID,
		Content:    ExtractString(data, "content"),
		SenderID:   ExtractString(data, "sender_id"),
		RoomID:     roomID,
		ReadStatus: readStatus,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}, nil
}

func parseReadStatus(data map[string]interface{}) (types.ReadStatus, error) {
	if rs, ok := data["read_status"].(string); ok {
		if rs == string(types.Read) || rs == string(types.Unread) {
			return types.ReadStatus(rs), nil
		}
	}
	return "", fmt.Errorf("Invalid read_status value")
}

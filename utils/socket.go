package utils

import (
	"fmt"
	"github.com/kwa0x2/Settle-Backend/domain"
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

func ExtractTime(data map[string]interface{}, key string) *time.Time {
	if value, ok := data[key].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, value) // rfc3339=2024-10-29T19:07:38.9537926Z
		if err == nil {
			return &parsedTime
		}
	}
	return nil
}

func ParseAttachment(data map[string]interface{}) (*domain.Attachment, error) {
	attachmentID, err := ParseObjectIDFromData(data, "ID")
	if err != nil {
		return nil, fmt.Errorf("invalid attachment ID format")
	}

	size, ok := data["Size"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid Size format")
	}

	createdAt := ExtractTime(data, "CreatedAt")
	if createdAt == nil {
		return nil, fmt.Errorf("invalid or missing CreatedAt field")
	}

	updatedAt := ExtractTime(data, "UpdatedAt")
	if updatedAt == nil {
		return nil, fmt.Errorf("invalid or missing UpdatedAt field")
	}

	return &domain.Attachment{
		ID:          attachmentID,
		Filename:    ExtractString(data, "Filename"),
		Size:        int64(size),
		Url:         ExtractString(data, "Url"),
		ContentType: ExtractString(data, "ContentType"),
		CreatedAt:   *createdAt,
		UpdatedAt:   *updatedAt,
	}, nil
}

func ParseUser(data map[string]interface{}) (*domain.User, error) {
	totalPlaytime, ok := data["TotalPlaytime"].(float64)
	if !ok {
		return nil, fmt.Errorf("Invalid TotalPlaytime format")
	}

	createdAt := ExtractTime(data, "CreatedAt")
	if createdAt == nil {
		return nil, fmt.Errorf("Invalid or missing CreatedAt field")
	}

	updatedAt := ExtractTime(data, "UpdatedAt")
	if updatedAt == nil {
		return nil, fmt.Errorf("Invalid or missing UpdatedAt field")
	}

	return &domain.User{
		ID:            ExtractString(data, "ID"),
		Name:          ExtractString(data, "Name"),
		Avatar:        ExtractString(data, "Avatar"),
		ProfileURL:    ExtractString(data, "ProfileURL"),
		TotalPlaytime: int(totalPlaytime),
		CreatedAt:     *createdAt,
		UpdatedAt:     *updatedAt,
	}, nil
}

func ParseRepliedMessage(data map[string]interface{}) (*domain.Message, error) {
	repliedMessageID, err := ParseObjectIDFromData(data, "ID")
	if err != nil {
		return nil, fmt.Errorf("Invalid replied message ID format")
	}

	roomID, err := ParseObjectIDFromData(data, "RoomID")
	if err != nil {
		return nil, fmt.Errorf("Invalid replied room ID format")
	}

	sender, err := ParseUser(data["Sender"].(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("Invalid sender format22")
	}

	readStatus, ok := data["ReadStatus"].(float64)
	if !ok {
		return nil, fmt.Errorf("Invalid ReadStatus format2")
	}

	createdAt := ExtractTime(data, "CreatedAt")
	updatedAt := ExtractTime(data, "UpdatedAt")

	return &domain.Message{
		ID:         repliedMessageID,
		Content:    ExtractString(data, "Content"),
		Sender:     sender,
		RoomID:     roomID,
		ReadStatus: int(readStatus),
		CreatedAt:  *createdAt,
		UpdatedAt:  *updatedAt,
	}, nil
}

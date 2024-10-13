package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserRoom struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"room_id" validate:"required"` // uuid.UUID
	UserID    string             `bson:"user_id" validate:"required"`
	Visible   bool               `bson:"visible" validate:"required"`
	CreatedAt time.Time          `bson:"created_at" `
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"` // Silinme tarihi
}

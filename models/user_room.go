package models

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserRoom struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"room_id" validate:"required"` // uuid.UUID
	UserID    string             `bson:"user_id" validate:"required"`
	Visible   bool               `bson:"visible" validate:"required"`
	CreatedAt time.Time          `bson:"created_at"  validate:"required"`
	UpdatedAt time.Time          `bson:"updated_at"  validate:"required"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"` // Silinme tarihi
}

func (u *UserRoom) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

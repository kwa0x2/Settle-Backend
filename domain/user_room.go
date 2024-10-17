package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	CollectionUserRoom = "user_rooms"
)

type UserRoom struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	RoomID    string        `bson:"room_id" validate:"required"` // uuid.UUID
	UserID    string        `bson:"user_id" validate:"required"`
	Visible   bool          `bson:"visible" validate:"required"`
	CreatedAt time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"` // Silinme tarihi
}

func (u *UserRoom) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserRoomRepository interface {
	Create(ctx context.Context, userRoom *UserRoom) (*mongo.InsertOneResult, error)
}

type UserRoomUsecase interface {
	Create(userRoom *UserRoom) error
}

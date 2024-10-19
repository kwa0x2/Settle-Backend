package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

const (
	CollectionRoom = "rooms"
)

type Room struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	CreatedUserID string        `bson:"created_user_id"   validate:"required"`
	LastMessageID bson.ObjectID `bson:"last_message_id"   validate:"required"`
	CreatedAt     time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt     time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt     *time.Time    `bson:"deleted_at,omitempty"`

	LastMessage *Message `bson:"last_message,omitempty"`
}

func (r *Room) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type RoomUsecase interface {
	FindAll() ([]Room, error)
}

type RoomRepository interface {
	UpdateByID(ctx context.Context, roomID bson.ObjectID, update bson.D) error
	Find(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]Room, error)
}

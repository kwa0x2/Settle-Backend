package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	CollectionMessage = "messages"
)

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Content        string             `bson:"content,omitempty"`
	SenderID       string             `bson:"sender_id" validate:"required"` // uuid.UUID
	RoomID         string             `bson:"room_id" validate:"required"`   // uuid.UUID
	RepliedMessage *Message           `bson:"replied_message,omitempty"`     // Referans verilen mesaj
	Attachment     *Attachment        `bson:"attachment,omitempty"`          // Eklenen dosya
	ReadStatus     types.ReadStatus   `bson:"read_status" validate:"required"`
	CreatedAt      time.Time          `bson:"created_at"  validate:"required"`
	UpdatedAt      time.Time          `bson:"updated_at"  validate:"required"`
	DeletedAt      *time.Time         `bson:"deleted_at,omitempty"`
}

func (m *Message) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type MessageRepository interface {
	Create(ctx context.Context, message *Message) (*mongo.InsertOneResult, error)
	UpdateByID(ctx context.Context, messageID bson.ObjectID, update bson.D) (interface{}, error)
}

type MessageUsecase interface {
	Create(message *Message) error
	SoftDelete(messageID bson.ObjectID) error
}

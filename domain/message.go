package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

const (
	CollectionMessage = "messages"
)

type Message struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	Content        string        `bson:"content,omitempty"`
	Sender         *User         `bson:"sender" validate:"required"`
	RoomID         bson.ObjectID `bson:"room_id" validate:"required"`
	RepliedMessage *Message      `bson:"replied_message,omitempty"` // Referans verilen mesaj
	Attachment     *Attachment   `bson:"attachment,omitempty"`      // Eklenen dosya
	ReadStatus     int           `bson:"read_status"`
	CreatedAt      time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt      time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt      *time.Time    `bson:"deleted_at,omitempty"`
}

func (m *Message) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type MessageRepository interface {
	Create(ctx context.Context, message *Message) (*mongo.InsertOneResult, error)
	UpdateByID(ctx context.Context, messageID bson.ObjectID, update bson.D) error
	Find(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]Message, error)
	FindOne(ctx context.Context, filter bson.D) (Message, error)
	GetDatabase() *mongo.Database
}

type MessageUsecase interface {
	CreateAndUpdateRoom(message *Message) error
	SoftDelete(messageID bson.ObjectID) error
	GetByRoomID(roomID bson.ObjectID, limit, offset int64) ([]Message, error)
	EditMessage(messageID bson.ObjectID, content string) error
}

type MessageHistoryRequest struct {
	RoomID bson.ObjectID `json:"RoomID"`
	Limit  int64         `json:"Limit"`
	Offset int64         `json:"Offset"`
}

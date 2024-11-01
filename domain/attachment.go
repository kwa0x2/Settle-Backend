package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	CollectionAttachment = "attachments"
)

type Attachment struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Filename    string        `bson:"filename"   validate:"required"`
	Size        int64         `bson:"size"   validate:"required"`
	Url         string        `bson:"url"   validate:"required"`
	ContentType string        `bson:"content_type"  validate:"required"`
	CreatedAt   time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt   time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt   *time.Time    `bson:"deleted_at,omitempty"`
}

func (a *Attachment) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

type AttachmentUsecase interface {
	Create(attachment *Attachment) error
}

type AttachmentRepository interface {
	Create(ctx context.Context, message *Attachment) (*mongo.InsertOneResult, error)
}

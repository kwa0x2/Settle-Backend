package domain

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Attachment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Filename    string             `bson:"filename"   validate:"required"`
	Size        int                `bson:"size"   validate:"required"`
	Url         string             `bson:"url"   validate:"required"`
	ContentType string             `bson:"content_type"  validate:"required"`
	CreatedAt   time.Time          `bson:"created_at"  validate:"required"`
	UpdatedAt   time.Time          `bson:"updated_at"  validate:"required"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty"`
}

func (a *Attachment) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

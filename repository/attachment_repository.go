package repository

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type attachmentRepository struct {
	database   *mongo.Database
	collection string
}

func NewAttachmentRepository(db *mongo.Database, collection string) domain.AttachmentRepository {
	return &attachmentRepository{
		database:   db,
		collection: collection,
	}
}

func (ar *attachmentRepository) Create(ctx context.Context, attachment *domain.Attachment) (*mongo.InsertOneResult, error) {
	collection := ar.database.Collection(ar.collection)

	return collection.InsertOne(ctx, attachment)
}

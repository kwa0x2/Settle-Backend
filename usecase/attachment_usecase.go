package usecase

import (
	"context"
	"github.com/kwa0x2/Settle-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type attachmentUsecase struct {
	attachmentRepository domain.AttachmentRepository
}

func NewAttachmentUsecase(attachmentRepository domain.AttachmentRepository) domain.AttachmentUsecase {
	return &attachmentUsecase{
		attachmentRepository: attachmentRepository,
	}
}

func (au *attachmentUsecase) Create(attachment *domain.Attachment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	attachment.CreatedAt = time.Now().UTC()
	attachment.UpdatedAt = time.Now().UTC()
	if err := attachment.Validate(); err != nil {
		return err
	}
	result, err := au.attachmentRepository.Create(ctx, attachment)
	if err != nil {
		return err
	}

	attachment.ID = result.InsertedID.(bson.ObjectID)

	return nil
}

package delivery

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/utils"
	"net/http"
)

type AttachmentDelivery struct {
	AttachmentUsecase domain.AttachmentUsecase
	S3                *s3.Client
	Env               *bootstrap.Env
}

func (ad *AttachmentDelivery) Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Form File Error: " + err.Error()})
		return
	}
	defer file.Close()

	fileURL, UploadErr := utils.UploadFile(file, header, ad.Env, ad.S3)
	if UploadErr != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Error uploading file to S3 bucket"})
		return
	}

	attachment := &domain.Attachment{
		Filename:    header.Filename,
		Size:        header.Size,
		Url:         fileURL,
		ContentType: header.Header.Get("Content-Type"),
	}

	if UpdateErr := ad.AttachmentUsecase.Create(attachment); UpdateErr != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Error updating user photo"})
		return
	}

	ctx.JSON(http.StatusCreated, attachment)
}

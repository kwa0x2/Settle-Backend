package route

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/api/http/delivery"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/repository"
	"github.com/kwa0x2/Settle-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewAttachmentRoute(env *bootstrap.Env, db *mongo.Database, group *gin.RouterGroup, s3 *s3.Client) {
	ar := repository.NewAttachmentRepository(db, domain.CollectionAttachment)
	md := &delivery.AttachmentDelivery{
		AttachmentUsecase: usecase.NewAttachmentUsecase(ar),
		S3:                s3,
		Env:               env,
	}

	group.POST("attachment", md.Upload)
}

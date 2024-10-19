package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/domain"
	"net/http"
)

type MessageDelivery struct {
	MessageUsecase domain.MessageUsecase
}

func (md *MessageDelivery) GetMessageHistory(ctx *gin.Context) {
	var req domain.MessageHistoryRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "JSON Bind Error"})
		return
	}

	messageHistoryData, err := md.MessageUsecase.GetByRoomID(req.RoomID, req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"limit": req.Limit,
		"page":  req.Offset,
		"data":  messageHistoryData,
	})

}

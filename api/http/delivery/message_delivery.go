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

	// Bind JSON request body to the MessageHistoryBody struct.
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "JSON Bind Error"})
		return
	}

	// Retrieve message history data using the provided room ID.
	messageHistoryData, err := md.MessageUsecase.GetByRoomID(req.RoomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, messageHistoryData)

}

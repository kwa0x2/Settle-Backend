package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/domain"
	"net/http"
)

type RoomDelivery struct {
	RoomUsecase domain.RoomUsecase
}

func (ad *RoomDelivery) GetRooms(ctx *gin.Context) {
	roomsData, err := ad.RoomUsecase.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, roomsData)
}

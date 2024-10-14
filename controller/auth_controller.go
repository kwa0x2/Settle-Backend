package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwa0x2/Settle-Backend/config"
	"github.com/kwa0x2/Settle-Backend/models"
	"github.com/kwa0x2/Settle-Backend/service"
	"github.com/kwa0x2/Settle-Backend/types"
	"github.com/kwa0x2/Settle-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

type IAuthController interface {
	SteamLogin(ctx *gin.Context)
	SteamCallback(ctx *gin.Context)
}

type authController struct {
	UserService     service.IUserService
	UserRoomService service.IUserRoomService
}

func NewAuthController(userService service.IUserService, userRoomService service.IUserRoomService) IAuthController {
	return &authController{
		UserService:     userService,
		UserRoomService: userRoomService,
	}
}

func (ctrl *authController) SteamLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, config.GetSteamLoginURL())
}

func (ctrl *authController) SteamCallback(ctx *gin.Context) {
	identity := ctx.Query("openid.identity")
	if identity == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Identity not found!"})
		return
	}

	steamID, err := utils.ExtractSteamID(identity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract Steam ID: " + err.Error()})
		return
	}

	userInfo, userInfoErr := utils.GetUserInfo(steamID)
	if userInfoErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information: " + userInfoErr.Error()})
		return
	}

	totalPlayTime, playTimeErr := utils.GetTotalPlaytime(steamID)
	if playTimeErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve playtime information: " + playTimeErr.Error()})
		return
	}

	if totalPlayTime < 30000 { // less than 500 hours
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Playtime must be more than 500 hours"})
		return
	}

	newUser := &models.User{
		ID:            userInfo.ID,
		Name:          userInfo.Name,
		Avatar:        userInfo.Avatar,
		ProfileURL:    userInfo.ProfileURL,
		TotalPlaytime: totalPlayTime,
		Role:          types.User,
	}

	roomId, uuidErr := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if uuidErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UUID format"})
		return
	}

	newUserRoom := &models.UserRoom{
		RoomID: roomId.String(),
		UserID: userInfo.ID,
	}

	err = ctrl.UserService.CreateAndJoinRoom(newUser, newUserRoom)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already registered"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while creating the user and joining the room: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_room": newUserRoom,
		"user":      newUser,
	})
}

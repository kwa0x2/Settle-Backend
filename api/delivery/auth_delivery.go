package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/models"
	"github.com/kwa0x2/Settle-Backend/types"
	"github.com/kwa0x2/Settle-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

type AuthController struct {
	UserUsecase domain.UserUsecase
	Env         *bootstrap.Env
}

func (ctrl *AuthController) SteamLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, utils.GetSteamLoginURL(ctrl.Env))
}

func (ac *AuthController) SteamCallback(ctx *gin.Context) {
	identity := ctx.Query("openid.identity")
	if identity == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Identity not found!"})
		return
	}

	steamID, err := utils.ExtractSteamID(identity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to extract Steam ID: " + err.Error()})
		return
	}

	userInfo, userInfoErr := utils.GetUserInfo(steamID, ac.Env)
	if userInfoErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to retrieve user information: " + userInfoErr.Error()})
		return
	}

	totalPlayTime, playTimeErr := utils.GetTotalPlaytime(steamID, ac.Env)
	if playTimeErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to retrieve playtime information: " + playTimeErr.Error()})
		return
	}

	if totalPlayTime < 30000 { // less than 500 hours
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Playtime must be more than 500 hours"})
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
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Invalid UUID format"})
		return
	}

	newUserRoom := &models.UserRoom{
		RoomID: roomId.String(),
		UserID: userInfo.ID,
	}

	err = ac.UserUsecase.CreateAndJoinRoom(newUser, newUserRoom)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "User already registered"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "An error occurred while creating the user and joining the room: " + err.Error()})
		return
	}

	accessToken, accessTokenErr := utils.CreateAccessToken(userInfo.ID, ac.Env.AccessTokenSecret, ac.Env.AccessTokenExpiryHour)
	if accessTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: accessTokenErr.Error()})
		return
	}

	refreshToken, refreshTokenErr := utils.CreateRefreshToken(userInfo.ID, ac.Env.RefreshTokenSecret, ac.Env.RefreshTokenExpiryHour)
	if refreshTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: refreshTokenErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

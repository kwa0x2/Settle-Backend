package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"github.com/kwa0x2/Settle-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

type AuthDelivery struct {
	UserUsecase domain.UserUsecase
	Env         *bootstrap.Env
}

func (ad *AuthDelivery) Login(ctx *gin.Context) {
	var req domain.LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "JSON Bind Error"})
		return
	}

	userInfo, userInfoErr := utils.GetUserInfo(req.SteamID, ad.Env.SteamApiKey)
	if userInfoErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to retrieve user information: " + userInfoErr.Error()})
		return
	}

	totalPlayTime, playTimeErr := utils.GetTotalPlaytime(req.SteamID, ad.Env.SteamApiKey)
	if playTimeErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to retrieve playtime information: " + playTimeErr.Error()})
		return
	}

	if totalPlayTime < 30000 { // less than 500 hours
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Playtime must be more than 500 hours"})
		return
	}

	newUser := &domain.User{
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

	newUserRoom := &domain.UserRoom{
		RoomID: roomId.String(),
		UserID: userInfo.ID,
	}

	accessToken, accessTokenErr := utils.CreateAccessToken(userInfo.ID, ad.Env.AccessTokenSecret, ad.Env.AccessTokenExpiryHour)
	if accessTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: accessTokenErr.Error()})
		return
	}

	refreshToken, refreshTokenErr := utils.CreateRefreshToken(userInfo.ID, ad.Env.RefreshTokenSecret, ad.Env.RefreshTokenExpiryHour)
	if refreshTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: refreshTokenErr.Error()})
		return
	}

	err := ad.UserUsecase.CreateAndJoinRoom(newUser, newUserRoom)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusOK, domain.LoginResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "An error occurred while creating the user and joining the room: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ad *AuthDelivery) RefreshToken(ctx *gin.Context) {
	var req domain.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid request body: " + err.Error()})
		return
	}

	userID, err := utils.IsAuthorized(req.RefreshToken, ad.Env.RefreshTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid refresh token: " + err.Error()})
		return
	}

	accessToken, accessTokenErr := utils.CreateAccessToken(userID, ad.Env.AccessTokenSecret, ad.Env.AccessTokenExpiryHour)
	if accessTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: accessTokenErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.RefreshResponse{
		AccessToken: accessToken,
	})

}

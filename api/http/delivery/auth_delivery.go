package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwa0x2/Settle-Backend/bootstrap"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/domain/types"
	"github.com/kwa0x2/Settle-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"time"
)

type AuthDelivery struct {
	UserUsecase domain.UserUsecase
	Env         *bootstrap.Env
}

func (ad *AuthDelivery) SteamLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, utils.GetSteamLoginURL(ad.Env.SteamRedirectUrl))
}

func (ad *AuthDelivery) SteamCallback(ctx *gin.Context) {
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

	userInfo, userInfoErr := utils.GetUserInfo(steamID, ad.Env.SteamApiKey)
	if userInfoErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to retrieve user information: " + userInfoErr.Error()})
		return
	}

	totalPlayTime, playTimeErr := utils.GetTotalPlaytime(steamID, ad.Env.SteamApiKey)
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

	accessToken, accessTokenErr := utils.CreateAccessToken(newUser, ad.Env.AccessTokenSecret, ad.Env.AccessTokenExpiryHour)
	if accessTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: accessTokenErr.Error()})
		return
	}

	refreshToken, refreshTokenErr := utils.CreateRefreshToken(newUser, ad.Env.RefreshTokenSecret, ad.Env.RefreshTokenExpiryHour)
	if refreshTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: refreshTokenErr.Error()})
		return
	}

	redirectUrl := fmt.Sprintf("http://100.64.75.37:4724/en/auth/login?access=%s&refresh=%s", accessToken, refreshToken)

	fmt.Println(redirectUrl)
	err = ad.UserUsecase.CreateAndJoinRoom(newUser, newUserRoom)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			time.Sleep(1000 * time.Millisecond)
			ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
			return
		}
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "An error occurred while creating the user and joining the room: " + err.Error()})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (ad *AuthDelivery) RefreshToken(ctx *gin.Context) {
	var req domain.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid request body: " + err.Error()})
		return
	}

	userData, err := utils.IsAuthorized(req.RefreshToken, ad.Env.RefreshTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid refresh token: " + err.Error()})
		return
	}

	accessToken, accessTokenErr := utils.CreateAccessToken(userData, ad.Env.AccessTokenSecret, ad.Env.AccessTokenExpiryHour)
	if accessTokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: accessTokenErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.RefreshResponse{
		AccessToken: accessToken,
	})
}

func (ad *AuthDelivery) CheckAuth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, true)
}

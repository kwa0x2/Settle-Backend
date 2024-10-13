package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/config"
	"github.com/kwa0x2/Settle-Backend/models"
	"github.com/kwa0x2/Settle-Backend/service"
	"github.com/kwa0x2/Settle-Backend/utils"
	"net/http"
)

type IAuthController interface {
	SteamLogin(ctx *gin.Context)
	SteamCallback(ctx *gin.Context)
}

type authController struct {
	UserService service.IUserService
}

func NewAuthController(userService service.IUserService) IAuthController {
	return &authController{
		UserService: userService,
	}
}

func (ctrl *authController) SteamLogin(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, config.GetSteamLoginURL())
}

func (ctrl *authController) SteamCallback(ctx *gin.Context) {
	identity := ctx.Query("openid.identity")
	if identity == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"erorr": "Identity not found!"})
		return
	}

	steamID, err := utils.ExtractSteamID(identity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Steam kimliği çıkarılamadı: " + err.Error()})
		return
	}

	userInfo, userInfoErr := utils.GetUserInfo(steamID)
	if userInfoErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı bilgileri alınamadı: " + userInfoErr.Error()})
		return
	}

	totalPlayTime, playTimeErr := utils.GetTotalPlaytime(steamID)
	if playTimeErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Oyun bilgileri alınamadı: " + playTimeErr.Error()})
		return
	}

	if totalPlayTime < 30000 { // 500 saaten az ise
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "oyun suresi 500satten az"})
		return
	}

	newUser := &models.User{
		ID:            userInfo.ID,
		Name:          userInfo.Name,
		Avatar:        userInfo.Avatar,
		ProfileURL:    userInfo.ProfileURL,
		TotalPlaytime: totalPlayTime,
	}

	if createErr := ctrl.UserService.Create(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı veritabanına eklenemedi: " + createErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newUser)
}

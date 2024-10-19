package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/Settle-Backend/domain"
	"github.com/kwa0x2/Settle-Backend/utils"
	"net/http"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Authorization header is missing"})
			ctx.Abort()
			return
		}

		user, err := utils.IsAuthorized(token, secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

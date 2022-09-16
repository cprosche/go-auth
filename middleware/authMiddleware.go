package middleware

import (
	"net/http"

	"github.com/cprosche/auth/utils"
	"github.com/gin-gonic/gin"
)

func ValidateTokenHandler(context *gin.Context) {
	// get the token from the header
	authCookie, err := context.Request.Cookie("Authorization")
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
	}
	signedToken := authCookie.Value[7:len(authCookie.Value)]

	// validate the token
	userId, err := utils.ValidateToken(signedToken)
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
	}

	// set user id and continue
	context.Set("userId", userId)
	context.Next()
}

package middleware

import (
	"net/http"

	"github.com/cprosche/auth/inits"
	"github.com/cprosche/auth/utils"
	"github.com/gin-gonic/gin"
)

func ValidateTokenHandler(context *gin.Context) {
	// get the token from the header
	var authCookie string
	var err error
	if inits.CURRENT_ENV == "development" {
		authCookie = context.GetHeader("Authorization")
	} else {
		authCookie, err = context.Cookie("Authorization")
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}

	signedToken := authCookie[7:]

	// validate the token
	userId, err := utils.ValidateToken(signedToken)
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
	}

	// set user id and continue
	context.Set("userId", userId)
	context.Next()
}

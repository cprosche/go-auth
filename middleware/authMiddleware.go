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
			context.Status(http.StatusUnauthorized)
			return
		}
	}

	signedToken := authCookie[7:]
	if signedToken == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	// validate the token
	userId, err := utils.ValidateToken(signedToken)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// set user id and continue
	context.Set("userId", userId)
	context.Next()
}

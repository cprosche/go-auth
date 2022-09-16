package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(context *gin.Context) {
	origin := context.Request.Header.Get("Origin")
	context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, accept, origin, Cache-Control, X-Requested-With, OPTIONS, Cache")
	context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if context.Request.Method == "OPTIONS" {
		context.Status(http.StatusNoContent)
		return
	}

	context.Next()
}

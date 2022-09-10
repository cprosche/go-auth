package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"

	ctrl "github.com/cprosche/auth/controllers"
)

// TODO: password validation
// TODO: ip address restriction? maybe on deployment server not app
// TODO: set up with pscale
// TODO: rate limiting?
// TODO: add signing and verifying with RSA key

func main() {
	loadEnv()
	// key := GetRSAKey()
	router := gin.Default()
	initRoutes(router)
	router.Run(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
}

func initRoutes(router *gin.Engine) {
	router.Use(CORSMiddleware)
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", ctrl.CreateNewUser) // create user
			auth.POST("/login", ctrl.LoginUser)
		}
		users := v1.Group("/users")
		{
			users.GET("/", ctrl.GetAllUsers)                          // get all users
			users.GET("/me", ctrl.ValidateTokenHandler, ctrl.GetUser) // get single user
			users.POST("/me", ctrl.UpdateUser)                        // update single user
			users.DELETE("/me", ctrl.DeleteUser)                      // delete single user
		}
	}
}

func loadEnv() {
	err := env.Load(".env")
	if err != nil {
		panic("Error loading environment variables")
	}
}

func CORSMiddleware(context *gin.Context) {
	origin := context.Request.Header.Get("Origin")
	context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, accept, origin, Cache-Control, X-Requested-With, OPTIONS, Cache")
	context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if context.Request.Method == "OPTIONS" {
		context.AbortWithStatus(204)
		return
	}

	context.Next()
}

func GetRSAKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(os.Getenv("PRIVATE_KEY")))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)
	return key
}

// nodemon --exec go run main.go --ext go

package main

import (
	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"

	fn "github.com/cprosche/auth/functions"
)

// nodemon --exec go run main.go --ext go

// TODO: password validation
// TODO: ip address restriction? maybe on deployment server not app
// TODO: set up with pscale

func main() {
	loadEnv()
	router := gin.Default()
	initRoutes(router)
	router.Run("localhost:8080")
}

func initRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", fn.CreateNewUser) // create user
			auth.POST("/login", fn.LoginUser)
		}
		users := v1.Group("/users")
		{
			users.GET("/", fn.GetAllUsers)     // get all users
			users.GET("/me", fn.GetUser)       // get single user
			users.POST("/me", fn.UpdateUser)   // update single user
			users.DELETE("/me", fn.DeleteUser) // update single user
		}
	}
}

func loadEnv() {
	err := env.Load(".env")
	if err != nil {
		panic("Error loading environment variables")
	}
}

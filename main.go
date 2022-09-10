package main

import (
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

func main() {
	loadEnv()
	router := gin.Default()
	initRoutes(router)
	router.Run(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
}

func initRoutes(router *gin.Engine) {
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
			users.DELETE("/me", ctrl.DeleteUser)                      // update single user
		}
	}
}

func loadEnv() {
	err := env.Load(".env")
	if err != nil {
		panic("Error loading environment variables")
	}
}

// nodemon --exec go run main.go --ext go

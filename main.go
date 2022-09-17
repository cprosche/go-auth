package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	ctrl "github.com/cprosche/auth/controllers"
	"github.com/cprosche/auth/inits"
	mw "github.com/cprosche/auth/middleware"
	"github.com/cprosche/auth/utils"
)

// TODO: password validation
// TODO: ip address restriction? maybe on deployment server not app
// TODO: set up with pscale
// TODO: rate limiting?

func init() {
	inits.LoadEnv()

	inits.PORT = os.Getenv("PORT")
	if inits.PORT == "" {
		panic("Port loading error")
	}

	inits.PRIVATE_KEY = os.Getenv("PRIVATE_KEY")
	if inits.PRIVATE_KEY == "" {
		panic("Private key loading error")
	}

	inits.RSA_KEY = utils.GetRSAKey(inits.PRIVATE_KEY)

	inits.AUTH_DSN = os.Getenv("AUTH_DSN")
	if inits.AUTH_DSN == "" {
		panic("Auth dsn loading error")
	}

	inits.CURRENT_ENV = os.Getenv("CURRENT_ENV")
}

func main() {
	router := gin.Default()
	initRoutes(router)
	router.Run(fmt.Sprintf("localhost:%s", inits.PORT))
}

func initRoutes(router *gin.Engine) {
	router.Use(mw.CORSMiddleware)
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", ctrl.CreateNewUser) // create user
			auth.POST("/login", ctrl.LoginUser)
		}
		users := v1.Group("/users")
		{
			users.GET("/", ctrl.GetAllUsers)                        // get all users
			users.GET("/me", mw.ValidateTokenHandler, ctrl.GetUser) // get single user
			users.POST("/me", ctrl.UpdateUser)                      // update single user
			users.DELETE("/me", ctrl.DeleteUser)                    // delete single user
		}
	}
}

// nodemon --exec go run main.go --ext go

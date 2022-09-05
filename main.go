package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Pw       string `json:"pw"`
}

func main() {
	loadEnv()
	router := gin.Default()
	router.GET("/users", getAllUsers)
	router.POST("/register", getAllUsers)
	fmt.Println("Listening at http://localhost:8080")
	router.Run("localhost:8080")
}

func loadEnv() {
	err := env.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func getAllUsers(context *gin.Context) {
	db := getAuthDbConnection()
	defer db.Close()

	results, err := db.Query("SELECT id, username, email, pw FROM users")
	if err != nil {
		panic(err.Error())
	}

	var users []User
	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.Username, &user.Email, &user.Pw)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}

	context.IndentedJSON(http.StatusAccepted, users)
}

func getAuthDbConnection() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

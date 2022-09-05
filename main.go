package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// nodemon --exec go run main.go --ext go

// TODO: login route
// TODO: password validation
// TODO: ip address restriction? maybe on deployment server not app
// TODO: jwt token return in login route (edcsa public private key pair?)
// TODO: set up with pscale

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
	router.POST("/register", createNewUser)

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

func createNewUser(context *gin.Context) {
	// get the new user from the request body
	var newUser User
	err := context.BindJSON(&newUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	// generate a hash from the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Pw), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "password hashing failed"})
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// insert user into db
	sql := "INSERT INTO users(username, email, pw) VALUES (?, ?, ?)"
	_, err = db.Exec(sql, newUser.Username, newUser.Email, hashedPassword)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	// query new user from db to return new id
	row := db.QueryRow("SELECT id, username, email, pw FROM users WHERE username = ?", newUser.Username)
	err = row.Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.Pw)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	// return new user id to client
	context.IndentedJSON(http.StatusAccepted, newUser.ID)
}

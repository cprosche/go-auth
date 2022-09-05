package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Pw        string `json:"pw"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func main() {
	loadEnv()
	router := gin.Default()

	router.GET("/users", getAllUsers)
	router.POST("/register", createNewUser)
	router.POST("/login", loginUser)

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
	context.Status(http.StatusAccepted)
}

func loginUser(context *gin.Context) {
	// get user from request
	var userFromRequest User
	err := context.BindJSON(&userFromRequest)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// lookup user in db
	var userFromDb User
	row := db.QueryRow("SELECT id, username, email, pw FROM users WHERE username = ? OR email = ?", userFromRequest.Username, userFromRequest.Email)
	err = row.Scan(&userFromDb.ID, &userFromDb.Username, &userFromDb.Email, &userFromDb.Pw)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	// verify the password from request and db match
	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Pw), []byte(userFromRequest.Pw))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// create jwt
	// TODO: add correct claims
	unsignedJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// sign jwt
	signedJwt, err := unsignedJwt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// attach a jwt token to reponse, security reference: https://dev.to/gkoniaris/how-to-securely-store-jwt-tokens-51cf
	context.SetCookie(
		"Authorization",
		fmt.Sprintf("Bearer %s", signedJwt),
		60*60*24,
		"/",
		"localhost",
		false,
		true)
	context.SetSameSite(http.SameSiteStrictMode)

	// return an ok status
	context.IndentedJSON(http.StatusOK, gin.H{"jwt": signedJwt})
}

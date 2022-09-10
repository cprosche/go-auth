package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type FullUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Pw        string `json:"pw"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type SafeUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func getAuthDbConnection() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

func ValidateToken(signedToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token expired")
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ValidateTokenHandler(context *gin.Context) {
	// get the token from the header
	bearer := context.GetHeader("Authorization")
	if bearer == "" {
		context.Status(http.StatusUnauthorized)
		return
	}
	signedToken := strings.Replace(bearer, "Bearer ", "", -1)

	// validate the token
	userId, err := ValidateToken(signedToken)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}
	context.Set("userId", userId)
	context.Next()
}

func GetUser(context *gin.Context) {
	// get user id from the moddleware
	userId, ok := context.Get("userId")
	if !ok {
		context.Status(http.StatusUnauthorized)
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// get user from db
	var user SafeUser
	row := db.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", userId)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	context.IndentedJSON(http.StatusOK, user)
}

func UpdateUser(context *gin.Context) {}

func DeleteUser(context *gin.Context) {}

func GetAllUsers(context *gin.Context) {
	db := getAuthDbConnection()
	defer db.Close()

	results, err := db.Query("SELECT id, username, email, pw, created_at, updated_at FROM users")
	if err != nil {
		panic(err.Error())
	}

	var users []FullUser
	for results.Next() {
		var user FullUser
		err = results.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Pw,
			&user.CreatedAt,
			&user.UpdatedAt)

		if err != nil {
			panic(err.Error())
		}

		users = append(users, user)
	}

	context.IndentedJSON(http.StatusAccepted, users)
}

func CreateNewUser(context *gin.Context) {
	// get the new user from the request body
	var newUser FullUser
	err := context.BindJSON(&newUser)
	if err != nil {
		context.Status(http.StatusBadGateway)
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
	context.Status(http.StatusCreated)
}

func LoginUser(context *gin.Context) {
	// get user from request
	var userFromRequest FullUser
	err := context.BindJSON(&userFromRequest)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// lookup user in db
	var userFromDb FullUser
	row := db.QueryRow("SELECT id, username, email, pw FROM users WHERE username = ? OR email = ?", userFromRequest.Username, userFromRequest.Email)
	err = row.Scan(&userFromDb.ID, &userFromDb.Username, &userFromDb.Email, &userFromDb.Pw)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	// verify the password from request and db match
	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Pw), []byte(userFromRequest.Pw))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// create jwt
	unsignedJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    "github.com/cprosche",
		Subject:   strconv.Itoa(userFromDb.ID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(4 * time.Hour).Unix(),
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
		true,
		true)
	context.SetSameSite(http.SameSiteStrictMode)

	// return token in body when developing
	if os.Getenv("CURRENT_ENV") == "development" {
		context.IndentedJSON(http.StatusOK, gin.H{"token": signedJwt})
		return
	}

	// return an ok status
	context.Status(http.StatusOK)
}

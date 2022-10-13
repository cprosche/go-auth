package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cprosche/auth/inits"
	"github.com/cprosche/auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type FullUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Pw        string `json:"pw"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type SafeUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func getAuthDbConnection() *sql.DB {
	db, err := sql.Open("mysql", inits.AUTH_DSN)
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

func GetUser(context *gin.Context) {
	// get user id from the middleware
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
	row := db.QueryRow("SELECT id, email, created_at, updated_at FROM users WHERE id = ?", userId)
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// TODO: add check on updated_at vs jwt issued, unauthorized if updated_at is after issued_at

	context.IndentedJSON(http.StatusOK, user)
}

// TODO: create update email route
// TODO: add verification
func UpdateUserEmail(context *gin.Context) {
	// get user id from the middleware
	userId, ok := context.Get("userId")
	if !ok {
		context.Status(http.StatusUnauthorized)
		return
	}

	// expected request format
	type RequestContract struct {
		Email string `json:"email"`
	}

	// get email from request
	var request RequestContract
	err := context.BindJSON(&request)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// update email in db
	_, err = db.Exec("UPDATE users SET email = ? WHERE id = ?", request.Email, userId)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// returned status of accepted for update
	context.Status(http.StatusOK)
}

// TODO: create update password route
func UpdateUserPassword(context *gin.Context) {
	// get user id from the middleware
	userId, ok := context.Get("userId")
	if !ok {
		context.Status(http.StatusUnauthorized)
		return
	}

	// expected request format
	type RequestContract struct {
		Pw string `json:"pw"`
	}

	// get email from request
	var request RequestContract
	err := context.BindJSON(&request)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// check for password validity
	if !utils.IsPasswordValid(request.Pw) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	// generate a hash from the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Pw), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "hashing error"})
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// update email in db
	_, err = db.Exec("UPDATE users SET pw = ? WHERE id = ?", hashedPassword, userId)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// returned status of accepted for update
	context.Status(http.StatusOK)
}

func DeleteUser(context *gin.Context) {
	// get user id from the middleware
	userId, ok := context.Get("userId")
	if !ok {
		context.Status(http.StatusUnauthorized)
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// delete user from db
	var user SafeUser
	row := db.QueryRow("SELECT id, email, created_at, updated_at FROM users WHERE id = ?", userId)
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		context.Status(http.StatusNoContent)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// TODO: add check on updated_at vs jwt issued, unauthorized if updated_at is after issued_at

	context.Status(http.StatusOK)
}

func GetAllUsers(context *gin.Context) {
	db := getAuthDbConnection()
	defer db.Close()

	results, err := db.Query("SELECT id, email, pw, created_at, updated_at FROM users")
	if err != nil {
		panic(err.Error())
	}

	var users []FullUser
	for results.Next() {
		var user FullUser
		err = results.Scan(
			&user.ID,
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
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "binding error"})
		return
	}

	// check for password validity
	if !utils.IsPasswordValid(newUser.Pw) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	// generate a hash from the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Pw), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "hashing error"})
		return
	}

	// connect to db
	db := getAuthDbConnection()
	defer db.Close()

	// insert user into db
	sql := "INSERT INTO users(email, pw) VALUES (?, ?)"
	_, err = db.Exec(sql, newUser.Email, hashedPassword)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "insert error"})
		return
	}

	// query new user from db to return new id
	row := db.QueryRow("SELECT id, email, pw FROM users WHERE email = ?", newUser.Email)
	err = row.Scan(&newUser.ID, &newUser.Email, &newUser.Pw)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "select error"})
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
	row := db.QueryRow("SELECT id, email, pw FROM users WHERE email = ?", userFromRequest.Email)
	err = row.Scan(&userFromDb.ID, &userFromDb.Email, &userFromDb.Pw)
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
	// TODO: decrease expiration length
	unsignedJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Issuer:    "github.com/cprosche",
		Subject:   strconv.Itoa(userFromDb.ID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(4 * time.Hour).Unix(),
	})

	// sign jwt
	signedJwt, err := unsignedJwt.SignedString(inits.RSA_KEY)
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}

	// attach a jwt token to reponse, security reference: https://dev.to/gkoniaris/how-to-securely-store-jwt-tokens-51cf
	context.SetCookie(
		"Authorization",
		fmt.Sprintf("Bearer %s", signedJwt),
		60*60*24,
		"",
		"localhost",
		true,
		true)
	context.SetSameSite(http.SameSiteLaxMode)

	// return an ok status
	context.Status(http.StatusOK)
}

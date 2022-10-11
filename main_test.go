package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := InitRoutes()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "pong", response.Body.String())
}

// TODO: register test
// TODO: login test
// TODO: get user test
// TODO: update user test
// TODO: delete user test

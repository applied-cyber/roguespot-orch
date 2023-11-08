package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AP struct {
	SSID string `json:"ssid"`
	MAC  string `json:"mac"`
}

func postAP(c *gin.Context) {
	var newAP AP

	if err := c.BindJSON(&newAP); err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, newAP)
}

func main() {
	router := gin.Default()
	router.POST("/log", postAP)
	router.Run()
}

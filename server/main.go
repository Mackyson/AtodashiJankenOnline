package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/:greet", func(c *gin.Context) {
		hello := c.Param("greet")
		c.JSON(http.StatusOK, gin.H{
			"greet": hello,
		})
	})
	router.Run(":" + port)
}

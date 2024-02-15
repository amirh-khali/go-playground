package main

import (
	"github.com/amirh-khali/go-playground/db"
	"github.com/amirh-khali/go-playground/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	db.Connect()
}

func main() {
	router := gin.Default()

	router.GET("/", homePage)

	recipesHandler := handler.NewRecipesHandler()
	router.GET("/recipes", recipesHandler.List)
	router.POST("/recipes", recipesHandler.Add)
	router.GET("/recipes/:id", recipesHandler.Get)
	router.PUT("/recipes/:id", recipesHandler.Update)
	router.DELETE("/recipes/:id", recipesHandler.Remove)

	err := router.Run()
	if err != nil {
		return
	}
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}

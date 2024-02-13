package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/", homePage)

	recipesHandler := NewRecipesHandler()
	router.GET("/recipes", recipesHandler.List)
	router.POST("/recipes", recipesHandler.Add)
	router.GET("/recipes/:id", recipesHandler.Get)
	router.PUT("/recipes/:id", recipesHandler.Update)
	router.DELETE("/recipes/:id", recipesHandler.Remove)

	_ = router.Run()
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}

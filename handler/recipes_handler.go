package handler

import (
	"github.com/amirh-khali/go-playground/db"
	"github.com/amirh-khali/go-playground/db/models"
	handlerModels "github.com/amirh-khali/go-playground/handler/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecipesHandler struct{}

func NewRecipesHandler() *RecipesHandler {
	return &RecipesHandler{}
}

func (h RecipesHandler) Add(c *gin.Context) {
	var createRecipeRequest handlerModels.CreateRecipeRequest
	if err := c.ShouldBindJSON(&createRecipeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe := models.Recipe{Name: createRecipeRequest.Name, Description: createRecipeRequest.Description}
	db.DB.Create(&recipe)

	c.JSON(http.StatusOK, gin.H{"data": recipe})
}

func (h RecipesHandler) List(c *gin.Context) {
	var allRecipes []models.Recipe
	db.DB.Find(&allRecipes)

	c.JSON(http.StatusOK, gin.H{"data": allRecipes})
}

func (h RecipesHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var recipe models.Recipe
	if err := db.DB.Where("id = ?", id).First(&recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(200, gin.H{"data": recipe})
}

func (h RecipesHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var updateRecipeRequest handlerModels.UpdateRecipeRequest
	if err := c.ShouldBindJSON(&updateRecipeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var recipe models.Recipe
	if err := db.DB.Where("id = ?", id).First(&recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.DB.Model(&recipe).Updates(updateRecipeRequest)

	c.JSON(http.StatusOK, gin.H{"data": recipe})
}

func (h RecipesHandler) Remove(c *gin.Context) {
	id := c.Param("id")

	var recipe models.Recipe
	if err := db.DB.Where("id = ?", id).First(&recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.DB.Delete(&recipe)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

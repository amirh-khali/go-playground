package main

import (
	"errors"
	"net/http"

	"github.com/amirh-khali/orderbook/pkg/recipes"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type RecipesHandler struct {
	store recipes.MemStore
}

func NewRecipesHandler() *RecipesHandler {
	return &RecipesHandler{store: recipes.NewMemStore()}
}

func (h RecipesHandler) Add(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceID := slug.Make(recipe.Name)
	if err := h.store.Add(resourceID, recipe); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h RecipesHandler) List(c *gin.Context) {
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, r)
}

func (h RecipesHandler) Get(c *gin.Context) {
	id := c.Param("id")

	recipe, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(200, recipe)
}

func (h RecipesHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.Update(id, recipe); err != nil {
		if errors.Is(err, recipes.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h RecipesHandler) Remove(c *gin.Context) {
	id := c.Param("id")

	if err := h.store.Remove(id); err != nil {
		if errors.Is(err, recipes.NotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

package models

type UpdateRecipeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/amirh-khali/orderbook/pkg/recipes"
	"github.com/gosimple/slug"
)

func main() {
	mux := http.NewServeMux()

	home := HomeHandler{}
	mux.Handle("/", &home)

	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)
	mux.Handle("/recipes", recipesHandler)
	mux.Handle("/recipes/", recipesHandler)

	_ = http.ListenAndServe(":8080", mux)
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{store: s}
}

var (
	RecipeRe       = regexp.MustCompile(`^/recipes/*$`)
	RecipeReWithId = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.Add(w, r)
		return
	case r.Method == http.MethodGet && RecipeReWithId.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodPut && RecipeReWithId.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && RecipeReWithId.MatchString(r.URL.Path):
		h.Remove(w, r)
		return
	default:
		return
	}
}

func (h *RecipesHandler) Add(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := slug.Make(recipe.Name)
	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RecipesHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	recipe, err := h.store.Get(matches[1])
	if err != nil {
		if errors.Is(err, recipes.NotFoundErr) {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBytes)
}

func (h *RecipesHandler) List(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBytes)
}

func (h *RecipesHandler) Update(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(matches[1], recipe); err != nil {
		if errors.Is(err, recipes.NotFoundErr) {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RecipesHandler) Remove(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Remove(matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("This is my homepage"))
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("404 Not Found"))
}

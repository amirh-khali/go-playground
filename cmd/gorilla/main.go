package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amirh-khali/orderbook/pkg/recipes"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

func main() {
	router := mux.NewRouter()

	home := HomeHandler{}
	router.HandleFunc("/", home.ServeHTTP)

	s := router.PathPrefix("/recipes").Subrouter()
	NewRecipesHandler(s)

	_ = http.ListenAndServe(":8080", router)
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(router *mux.Router) *RecipesHandler {
	handler := &RecipesHandler{store: recipes.NewMemStore()}

	router.HandleFunc("", handler.List).Methods("GET")
	router.HandleFunc("/", handler.Add).Methods("POST")
	router.HandleFunc("/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/{id}", handler.Update).Methods("PUT")
	router.HandleFunc("/{id}", handler.Remove).Methods("DELETE")

	return handler
}

func (h RecipesHandler) Add(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w)
		return
	}

	resourceID := slug.Make(recipe.Name)
	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h RecipesHandler) List(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w)
		return
	}

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBytes)
}

func (h RecipesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	recipe, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, recipes.NotFoundErr) {
			NotFoundHandler(w)
			return
		}

		InternalServerErrorHandler(w)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBytes)
}

func (h RecipesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w)
		return
	}

	if err := h.store.Update(id, recipe); err != nil {
		if errors.Is(err, recipes.NotFoundErr) {
			NotFoundHandler(w)
			return
		}
		InternalServerErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h RecipesHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Remove(id); err != nil {
		InternalServerErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("This is my home page"))
}

func InternalServerErrorHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("404 Not Found"))
}

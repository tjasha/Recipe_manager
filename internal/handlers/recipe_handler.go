package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tjasha/Recipe_manager/internal/models"
)

// Load recipes from JSON file
func loadRecipes() ([]models.Recipe, error) {
	data, err := ioutil.ReadFile("data/recipes.json")
	if err != nil {
		return nil, err
	}

	var recipes []models.Recipe
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return nil, err
	}

	return recipes, nil
}

// Handler to serve recipes
func RecipesHandler(w http.ResponseWriter, r *http.Request) {
	recipes, err := loadRecipes()
	if err != nil {
		fmt.Println("failed to load:", err)
		http.Error(w, "Failed to load recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

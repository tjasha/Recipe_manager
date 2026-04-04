package service

import "github.com/tjasha/Recipe_manager/internal/model"

func CalculateIngredientsPerServing(recipe model.Recipe, portions int) ([]model.IngredientInRecipe, error) {

	// get list of ingredients
	ingredients := recipe.Ingredients
	// for loop to go over them
	for i := 0; i < len(ingredients); i++ {
		ingredients[i].Quantity = ingredients[i].Quantity / float64(recipe.Portion) * float64(portions)
	}

	return ingredients, nil
}

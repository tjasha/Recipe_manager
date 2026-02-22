package models

import "time"

// User struct
type User struct {
	ID          uint   `json:"id"`
	UserName    string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessLevel int    `json:"accessLevel"`
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

// Ingredient struct
type Ingredient struct {
}

// IngredientInRecipe struct
type IngredientInRecipe struct {
	ID         int64   `json:"id"`
	Ingredient string  `json:"ingredient"`
	Unit       string  `json:"unit"`     // kg, l, cup
	Quantity   float64 `json:"quantity"` // from other table
}

// RecipeInstructions struct
type RecipeInstructions struct {
	ID              int64  `json:"id"`
	StepSequence    int    `json:"stepSequence"`
	StepDescription string `json:"stepDescription"`
}

// Recipe struct
type Recipe struct {
	ID              int64                `json:"id"`
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	Portion         int                  `json:"portion"`
	PreparationTime int                  `json:"preparationTime"`
	CookingTime     int                  `json:"cookingTime"`
	CreatedAt       time.Time            `json:"createdAt"`
	ModifiedAt      time.Time            `json:"modifiedAt"`
	Author          User                 `json:"author"`
	Published       bool                 `json:"published"`
	Ingredients     []IngredientInRecipe `json:"ingredients"`
	Instructions    []RecipeInstructions `json:"instructions"`
}

package model

import (
	"time"
)

// User struct
type User struct {
	ID          uint      `json:"id"`
	UserName    string    `json:"name"`
	Email       string    `json:"email"`
	Password    *string   `json:"-"`
	AccessLevel int       `json:"accessLevel"`
	State       *string   `json:"state,omitempty"`
	GoogleID    *string   `json:"googleID,omitempty"`
	CreatedAt   time.Time `json:"-"`
	ModifiedAt  time.Time `json:"-"`
}

// Recipe struct
type Recipe struct {
	ID              int64                `json:"id"`
	Title           string               `json:"title"`
	Description     *string              `json:"description"`
	Portion         int                  `json:"portion"`
	PreparationTime *int                 `json:"preparationTime"`
	CookingTime     *int                 `json:"cookingTime"`
	Nutrition       *Nutrition           `json:"nutrition"`
	Author          *User                `json:"author"`
	Published       *bool                `json:"published"`
	ImageURL        *string              `json:"imageURL"`
	Ingredients     []IngredientInRecipe `json:"ingredients"`
	Instructions    []RecipeInstruction  `json:"instructions"`
	CreatedAt       time.Time            `json:"createdAt"`
	ModifiedAt      time.Time            `json:"modifiedAt"`
}

// Nutrition struct
type Nutrition struct {
	ID           int64    `json:"id"`
	Energy       *int     `json:"energy"`
	Calories     *int     `json:"calories"`
	Fat          *float64 `json:"fat"`
	SaturatedFat *float64 `json:"saturatedFat"`
	Sodium       *float64 `json:"sodium"`
	Fiber        *float64 `json:"fiber"`
	Carbohydrate *float64 `json:"carbohydrate"`
	Sugar        *float64 `json:"sugar"`
	Protein      *float64 `json:"protein"`
	Salt         *float64 `json:"salt"`
}

// Ingredient struct
type Ingredient struct {
	ID         int64   `json:"id"`
	Ingredient string  `json:"ingredient"`
	Unit       string  `json:"unit"` // kg, l, cup
	ImageURL   *string `json:"imageURL"`
}

// IngredientInRecipe struct
type IngredientInRecipe struct {
	ID         int64      `json:"id"`
	Ingredient Ingredient `json:"ingredient"`
	Quantity   float64    `json:"quantity"` // from other table
}

// RecipeInstruction struct
type RecipeInstruction struct {
	ID              int64  `json:"id"`
	StepSequence    int    `json:"stepSequence"`
	StepDescription string `json:"stepDescription"`
	ImageURL        string `json:"image"`
}

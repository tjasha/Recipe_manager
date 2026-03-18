package repository

import (
	"context"

	"github.com/tjasha/Recipe_manager/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllRecipes(ctx context.Context) ([]model.Recipe, error)
	GetRecipeByID(id int64) (model.Recipe, error)
	GetIngredientsByRecipeID(Context context.Context, id int64) ([]model.IngredientInRecipe, error)
	GetInstructionsByRecipeID(Context context.Context, id int64) ([]model.RecipeInstruction, error)
}

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tjasha/Recipe_manager/internal/model"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// CreateUser creates new user and returns its ID.
// It's gonna be used only for Goole login, it doesn't save password.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (user_name, email, google_id, access_level, created_at, modified_at) 
		 VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		 RETURNING id`,
		user.UserName, user.Email, user.GoogleID, user.AccessLevel,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByEmail retrieves a user from the database by their email address.
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_name, email, google_id, access_level, created_at, modified_at, state 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.UserName, &user.Email, &user.GoogleID, &user.AccessLevel, &user.CreatedAt, &user.ModifiedAt, &user.State)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

//func (r *PostgresRepository) AuthenticateUser(ctx context.Context, user *model.User) (int, string, int, error) {
//	var id, accessLevel int
//	var UserName string
//
//	err := r.pool.QueryRow(ctx,
//		"INSERT INTO users (user_name, email, password, access_level) VALUES ($1, $2, $3, $4) RETURNING id, user_name, access_level",
//		user.UserName, user.Email, user.Password, user.AccessLevel,
//	).Scan(&id, &UserName, &accessLevel)
//
//	return id, UserName, accessLevel, err
//}

// GetAllRecipes retrieves all recipes from the database.
func (r *PostgresRepository) GetAllRecipes(ctx context.Context) ([]model.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, portion, preparation_time, cooking_time, published, image_url,
			created_at, modified_at
		FROM recipe 
		ORDER BY modified_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []model.Recipe

	//scan every results
	for rows.Next() {
		var recipe model.Recipe
		err := rows.Scan(
			&recipe.ID,
			&recipe.Title,
			&recipe.Description,
			&recipe.Portion,
			&recipe.PreparationTime,
			&recipe.CookingTime,
			&recipe.Published,
			&recipe.ImageURL,
			&recipe.CreatedAt,
			&recipe.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

// GetRecipeByID retrieves a recipe from the database by its ID.
func (r *PostgresRepository) GetRecipeByID(id int64) (model.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, r.title, r.description, r.portion, r.preparation_time, r.cooking_time, r.published, 
       r.image_url, r.created_at, r.modified_at,
       u.id, u.user_name, 
       n.id, n.energy, n.calories, n.fat, n.saturated_fat, n.sodium, n.fiber, n.carbohydrate, n.sugar, 
       n.protein, n.salt
		from recipe r  
		INNER JOIN nutrition n ON r.nutrition_id=n.id
		INNER JOIN users u ON r.author=u.id 
		where r.id = $1
		`

	var recipe model.Recipe
	var nutrition model.Nutrition
	var author model.User

	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&recipe.ID,
		&recipe.Title,
		&recipe.Description,
		&recipe.Portion,
		&recipe.PreparationTime,
		&recipe.CookingTime,
		&recipe.Published,
		&recipe.ImageURL,
		&recipe.CreatedAt,
		&recipe.ModifiedAt,
		&author.ID,
		&author.UserName,
		nutrition.ID,
		&nutrition.Energy,
		&nutrition.Calories,
		&nutrition.Fat,
		&nutrition.SaturatedFat,
		&nutrition.Sodium,
		&nutrition.Fiber,
		&nutrition.Carbohydrate,
		&nutrition.Sugar,
		&nutrition.Protein,
		&nutrition.Salt,
	)
	if err != nil {
		return recipe, err
	}

	recipe.Author = &author
	recipe.Nutrition = &nutrition

	return recipe, nil
}

func (r *PostgresRepository) GetIngredientsByRecipeID(ctx context.Context, id int64) ([]model.IngredientInRecipe, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select i.id, i.ingredient, i.unit, i.image_url, 
       	ri.quantity, 
		from ingredient i 
		join recipe_ingredient ri on i.id = ri.ingredient_id
		where ri.recipe_id = $1
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredient model.Ingredient
	var ingredientsInRecipe []model.IngredientInRecipe

	//scan every results
	for rows.Next() {
		var ingredientInRecipe model.IngredientInRecipe
		err := rows.Scan(
			ingredient.ID,
			ingredient.Ingredient,
			ingredient.Unit,
			&ingredient.ImageURL,
			ingredientInRecipe.Quantity,
		)
		if err != nil {
			return nil, err
		}

		ingredientsInRecipe = append(ingredientsInRecipe, ingredientInRecipe)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ingredientsInRecipe, nil
}

func (r *PostgresRepository) GetInstructionsByRecipeID(ctx context.Context, id int64) ([]model.RecipeInstruction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, recipe_id, step_sequence, step_description, image_url 
		from instruction 
		where recipe_id = $1
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructions []model.RecipeInstruction

	//scan every results
	for rows.Next() {
		var instruction model.RecipeInstruction
		err := rows.Scan(
			instruction.ID,
			instruction.StepSequence,
			instruction.StepDescription,
			&instruction.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, instruction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instructions, nil
}

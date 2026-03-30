package repository

import (
	"context"
	"database/sql"
	"log"
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
func (r *PostgresRepository) GetAllPublishedRecipes(ctx context.Context) ([]model.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, portion, preparation_time, cooking_time, published, image_url,
			created_at, modified_at
		FROM recipe 
		WHERE published = true
		ORDER BY modified_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes = make([]model.Recipe, 0)

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
       n.id, n.calories, n.fat, n.sodium, n.fiber, n.carbohydrate, n.sugar, n.protein
		from recipe r  
		LEFT JOIN nutrition n ON r.id = n.recipe_id
		LEFT JOIN users u ON r.author=u.id 
		where r.id = $1
		`

	var recipe model.Recipe

	//nullable fields (they may not exist, which would cause an error while scanning
	var authorID sql.NullInt64
	var authorName sql.NullString
	var nutritionID sql.NullInt64
	var calories sql.NullFloat64
	var fat sql.NullFloat64
	var sodium sql.NullFloat64
	var fiber sql.NullFloat64
	var sugar sql.NullFloat64
	var protein sql.NullFloat64
	var carbohydrate sql.NullFloat64

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
		&authorID,
		&authorName,
		&nutritionID,
		&calories,
		&fat,
		&sodium,
		&fiber,
		&carbohydrate,
		&sugar,
		&protein,
	)
	if err != nil {
		log.Println("recipe scan", err)
		return recipe, err
	}

	//check that author exist
	if authorID.Valid {
		recipe.Author = &model.User{
			ID:       uint(authorID.Int64),
			UserName: authorName.String,
		}

	}

	//check that nutrition exist
	if nutritionID.Valid {
		caloriesVal := int(calories.Float64)
		recipe.Nutrition = &model.Nutrition{
			ID:           nutritionID.Int64,
			Calories:     &caloriesVal,
			Fat:          &fat.Float64,
			Sodium:       &sodium.Float64,
			Fiber:        &fiber.Float64,
			Carbohydrate: &carbohydrate.Float64,
			Sugar:        &sugar.Float64,
			Protein:      &protein.Float64,
		}
	}

	return recipe, nil
}

func (r *PostgresRepository) GetIngredientsByRecipeID(ctx context.Context, id int64) ([]model.IngredientInRecipe, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select i.id, i.ingredient, i.unit, i.image_url, 
       	ri.quantity, ri.id 
		from ingredient i 
		join recipe_ingredient ri on i.id = ri.ingredient_id
		where ri.recipe_id = $1
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredientsInRecipe []model.IngredientInRecipe

	//scan every results
	for rows.Next() {
		var ingredientInRecipe model.IngredientInRecipe
		err := rows.Scan(
			&ingredientInRecipe.Ingredient.ID,
			&ingredientInRecipe.Ingredient.Ingredient,
			&ingredientInRecipe.Ingredient.Unit,
			&ingredientInRecipe.Ingredient.ImageURL,
			&ingredientInRecipe.Quantity,
			&ingredientInRecipe.ID,
		)
		if err != nil {
			log.Println("ingredient scan err:", err)
			return nil, err
		}

		log.Println("ingredientInRecipe:", ingredientInRecipe)
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

	query := `select id, step_sequence, step_description, image_url 
		from instruction 
		where recipe_id = $1
		order by step_sequence ASC
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
			&instruction.ID,
			&instruction.StepSequence,
			&instruction.StepDescription,
			&instruction.ImageURL,
		)
		if err != nil {
			log.Println("instruction scan err:", err)
			return nil, err
		}

		instructions = append(instructions, instruction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instructions, nil
}

// GetAllRecipesFromUser retrieves all recipes from specific user from the database.
func (r *PostgresRepository) GetAllRecipesFromUser(ctx context.Context, userID uint) ([]model.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, portion, preparation_time, cooking_time, published, image_url,
			created_at, modified_at
		FROM recipe
		WHERE author = $1
		ORDER BY modified_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes = make([]model.Recipe, 0)

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

func (r *PostgresRepository) GetAllIngredients() ([]model.Ingredient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select i.id, i.ingredient, i.unit
		from ingredient i 
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []model.Ingredient

	//scan every results
	for rows.Next() {
		var ingredient model.Ingredient
		err := rows.Scan(
			&ingredient.ID,
			&ingredient.Ingredient,
			&ingredient.Unit,
		)
		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, ingredient)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *PostgresRepository) CreateRecipe(ctx context.Context, recipeData *model.RecipeForm) (*model.Recipe, error) {

	var recipe model.Recipe

	// we're using transaction here, we want to rollback if any of the queries isn't working properly
	// beginning of transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return &recipe, err
	}
	// makes sure that transaction is rolled back if any of the queries fail
	defer tx.Rollback(ctx)

	// get user id from session
	userId := ctx.Value("userID").(uint)

	// insert recipe and get id
	var recipeId int64
	err = tx.QueryRow(ctx, `
		INSERT INTO recipe (title, description, portion, preparation_time, cooking_time, 
		                    published, image_url, author)
		    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		    RETURNING id`,
		recipeData.Title, recipeData.Description, recipeData.Portion, recipeData.PreparationTime,
		recipeData.CookingTime, recipeData.Published, recipeData.ImageURL, userId,
	).Scan(&recipeId)
	if err != nil {
		return &recipe, err
	}

	// insert  nutrition
	_, err = tx.Exec(ctx, `
		INSERT INTO nutrition (calories, fat, sodium, fiber, carbohydrate, sugar, protein, recipe_id)
		    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		recipeData.Nutrition.Calories, recipeData.Nutrition.Fat, recipeData.Nutrition.Sodium,
		recipeData.Nutrition.Fiber, recipeData.Nutrition.Carbohydrate, recipeData.Nutrition.Sugar,
		recipeData.Nutrition.Protein, recipeId,
	)
	if err != nil {
		log.Println("nutrition scan err:", err)
		return &recipe, err
	}

	// insert ingredients
	for _, ingredient := range recipeData.Ingredients {
		_, err = tx.Exec(ctx, `
			INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity)
			VALUES ($1, $2, $3)`,
			recipeId, ingredient.IngredientId, ingredient.Quantity,
		)
		if err != nil {
			return &recipe, err
		}
	}

	// insert instructions
	for i, instruction := range recipeData.Instructions {
		_, err = tx.Exec(ctx, `	
			INSERT INTO instruction (recipe_id, step_sequence, step_description)
			    VALUES ($1, $2, $3)`,
			recipeId, i+1, instruction.StepDescription,
		)
		if err != nil {
			return &recipe, err
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return &recipe, err
	}

	recipe.ID = recipeId
	log.Println("Recipe created successfully", recipeId)
	return &recipe, nil
}

func (r *PostgresRepository) DeleteRecipe(ctx context.Context, recipeId int64, userId uint) error {

	query := ` DELETE FROM recipe 
		    WHERE id = $1
		    AND author = $2`

	ex, err := r.pool.Exec(ctx, query, recipeId, userId)
	if err != nil {
		return err
	}

	// check that something was deleted
	if ex.RowsAffected() == 0 {
		log.Println("No rows")
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresRepository) PublishRecipe(ctx context.Context, recipeId int64, newStatus bool, userId uint) error {

	query := `UPDATE recipe
			SET published = $1, modified_at = NOW()
			WHERE id = $2
			AND author = $3`

	ex, err := r.pool.Exec(ctx, query, newStatus, recipeId, userId)
	if err != nil {
		log.Println("ERROR:", err)
		return err
	}

	if ex.RowsAffected() == 0 {
		log.Println("No rows")
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresRepository) UpdateRecipe(ctx context.Context, recipeData *model.RecipeForm, userId uint) error {

	// we're using transaction here, we want to rollback if any of the queries isn't working properly
	// beginning of transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	// makes sure that transaction is rolled back if any of the queries fail
	defer tx.Rollback(ctx)

	// update recipe
	ex, err := tx.Exec(ctx, `
		UPDATE recipe
		SET title = $1, description = $2, portion = $3, preparation_time = $4, cooking_time = $5,
		    published = $6, image_url = $7, modified_at = NOW()
		WHERE id = $8
		AND author = $9`,
		recipeData.Title, recipeData.Description, recipeData.Portion, recipeData.PreparationTime,
		recipeData.CookingTime, recipeData.Published, recipeData.ImageURL, recipeData.ID, userId,
	)
	if err != nil {
		log.Println("ERROR recipe update:", err)
		return err
	}

	if ex.RowsAffected() == 0 {
		log.Println("No rows for recipe updated")
		return sql.ErrNoRows
	}

	// update  nutrition
	if recipeData.Nutrition == nil {
		_, err = tx.Exec(ctx, `DELETE FROM nutrition WHERE recipe_id = $1`, recipeData.ID)
		if err != nil {
			return err
		}
	} else {
		ex, err = tx.Exec(ctx, `
			UPDATE nutrition
			SET calories = $1, fat = $2, sodium = $3, fiber = $4, carbohydrate = $5, sugar = $6, protein = $7
			WHERE recipe_id = $8`,
			recipeData.Nutrition.Calories, recipeData.Nutrition.Fat, recipeData.Nutrition.Sodium,
			recipeData.Nutrition.Fiber, recipeData.Nutrition.Carbohydrate, recipeData.Nutrition.Sugar,
			recipeData.Nutrition.Protein, recipeData.ID,
		)
		if err != nil {
			return err
		}
		if ex.RowsAffected() == 0 {
			_, err = tx.Exec(ctx, `
				INSERT INTO nutrition (calories, fat, sodium, fiber, carbohydrate, sugar, protein, recipe_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				recipeData.Nutrition.Calories, recipeData.Nutrition.Fat, recipeData.Nutrition.Sodium,
				recipeData.Nutrition.Fiber, recipeData.Nutrition.Carbohydrate, recipeData.Nutrition.Sugar,
				recipeData.Nutrition.Protein, recipeData.ID,
			)
			if err != nil {
				log.Println("nutrition scan err:", err)
				return err
			}
		}
	}

	// update ingredients
	// delete all old ingredients for this recipe
	query := `DELETE FROM recipe_ingredient
		WHERE recipe_id = $1`
	ex, err = tx.Exec(ctx, query, recipeData.ID)
	if err != nil {
		return err
	}

	// insert new ingredients
	for _, ingredient := range recipeData.Ingredients {
		_, err = tx.Exec(ctx, `
			INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity)
			VALUES ($1, $2, $3)`,
			recipeData.ID, ingredient.IngredientId, ingredient.Quantity,
		)
		if err != nil {
			return err
		}
	}

	// update instructions
	// delete all old instructions for this recipe
	query = `DELETE FROM instruction
		WHERE recipe_id = $1`
	_, err = tx.Exec(ctx, query, recipeData.ID)
	if err != nil {
		return err
	}

	// insert new instructions
	for i, instruction := range recipeData.Instructions {
		ex, err = tx.Exec(ctx, `	
			INSERT INTO instruction (recipe_id, step_sequence, step_description)
			    VALUES ($1, $2, $3)`,
			recipeData.ID, i+1, instruction.StepDescription,
		)
		if err != nil {
			return err
		}
		if ex.RowsAffected() == 0 {
			log.Println("`rows for new instructions not added")
			return sql.ErrNoRows
		}
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	log.Println("Recipe updated successfully")
	return nil
}

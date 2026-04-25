package service

type RecipeFilters struct {
	AuthorID     *int64
	MinCalories  *int
	MaxCalories  *int
	MinProtein   *float64
	MaxProtein   *float64
	MaxTotalTime *int
	//IngredientIDs []int64
	//IngredientMode *string
}

type RecipeSort string

const (
	SortModifiedAtDesc RecipeSort = "modified_at_desc"
	SortModifiedAtAsc  RecipeSort = "modified_at_asc"
	SortCaloriesAsc    RecipeSort = "calories_asc"
	SortCaloriesDesc   RecipeSort = "calories_desc"
	SortTotalTimeAsc   RecipeSort = "total_time_asc"
	SortTotalTimeDesc  RecipeSort = "total_time_desc"
)

type Pagination struct {
	Limit  int
	Offset int
}

type ListRecipesParams struct {
	Filters    RecipeFilters
	Sort       RecipeSort
	Pagination Pagination
}

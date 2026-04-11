import { useState, useEffect } from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import api from '../api/client';
import { toast } from 'react-toastify';

interface IngredientInRecipe {
    id: number;
    ingredient: {
        id: number;
        ingredient: string;
        unit: string;
        imageURL: string | null;
    };
    quantity: number;
}

interface Instruction {
    id: number;
    stepSequence: number;
    stepDescription: string;
}

interface Nutrition {
    id: number;
    calories: number | null;
    fat: number | null;
    sodium: number | null;
    fiber: number | null;
    carbohydrate: number | null;
    sugar: number | null;
    protein: number | null;
}

interface FullRecipe {
    id: number;
    title: string;
    description: string | null;
    author: {
        id: number;
        name: string;
    } | null;
    imageURL: string | null;
    portion: number;
    preparationTime: number | null;
    cookingTime: number | null;
    ingredients: IngredientInRecipe[];
    instructions: Instruction[];
    nutrition: Nutrition | null;
    published: boolean | null;
}

export default function MyRecipePage() {
    const { id } = useParams<{ id: string }>(); // resd from URL
    const [recipe, setRecipe] = useState<FullRecipe | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();


    useEffect(() => {
        const fetchRecipe = async () => {
            try {
                const response = await api.get(`/myrecipe/${id}`);
                setRecipe(response.data);
            } catch (err: any) {
                //user is not logged in
                if (err.response && err.response.status === 401) {
                    setError('You must be logged in to view this page.');
                } else {
                    setError('Failed to fetch recipe.');
                }
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        
        fetchRecipe();
    }, [id]);

    //delete recipe
    const handleDelete = async (recipeId: number, recipeTitle: string) => {
        // confirm deletion
        if (!window.confirm(`Are you sure you want to delete "${recipeTitle}"?`)) {
            return;
        }

        try {
            const response = await api.delete(`/deleteRecipe/${recipeId}`);
            toast.success('Recipe deleted successfully!');
            navigate(`/myrecipes`);
        } catch (err: any) {
            if (err.response && err.response.status === 404) {
                toast.error("Error: This recipe was already deleted or you don't have permission.");
            } else if (err.response && err.response.status === 401) {
                toast.error("Error: You must be logged in to delete a recipe.");
            } else {
                toast.error(err.response?.data?.message || err.message || "An unknown error occurred.");
            }
        }
    };

    const handleTogglePublish = async (recipeId: number, currentStatus: boolean) => {
        const newStatus = !currentStatus;
        try {
             await api.patch(`/myrecipe/${recipeId}/publish`, {
                published: newStatus
            });

            if (recipe) {
                setRecipe({
                    ...recipe,
                    published: newStatus, //only update publish, keep other data
                });
            }

        } catch (err: any) {
            toast.error(`Error: ${err.response?.data?.message || err.message || 'Failed to update status.'}`);
        }
    };


    if (loading) return <div className="container mt-4"><h2>Loading...</h2></div>;
    if (error) return <div className="container mt-4"><div className="alert alert-danger">{error}</div></div>;
    if (!recipe) return <div className="container mt-4"><h2>Recipe not found.</h2></div>;

    return (
        <div className="container mt-5">
            <div className="row">
                <div className="col-md-8 offset-md-2">
                    <h1>{recipe.title}</h1>

                    <div className="d-grid gap-2 d-md-flex justify-content-md-center">
                        <button
                            className={`btn btn-sm me-2 ${recipe.published ? 'btn-outline-success' : 'btn-success'}`}
                            onClick={(e) => {
                                e.stopPropagation();
                                handleTogglePublish(recipe.id, recipe.published);
                            }}
                        >
                            {recipe.published === true ? 'Unpublish' : 'Publish'}
                        </button>
                        <button className="btn btn-primary me-2"
                        onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/myrecipe/${recipe.id}/edit`);
                        }}
                        >Edit</button>
                        <button className="btn btn-outline-secondary"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleDelete(recipe.id, recipe.title);
                                }}
                        >Delete</button>
                    </div>

                    <br/>


                        {recipe.imageURL &&
                            <img src={recipe.imageURL} alt={recipe.title} className="img-fluid rounded mb-4"/>}

                        <p className="lead">{recipe.description}</p>

                        <div className="my-4">
                            <span><strong>Servings:</strong> {recipe.portion} | </span>
                            <span><strong>Prep Time:</strong> {recipe.preparationTime} min | </span>
                            <span><strong>Cook Time:</strong> {recipe.cookingTime} min</span>
                        </div>

                        <div className="row">
                            <div className="col-md-5">
                                <h3>Ingredients</h3>
                                <ul className="list-group">
                                    {recipe.ingredients.map(ingInRecipe => (
                                        <li key={ingInRecipe.id} className="list-group-item d-flex align-items-center">

                                            {/* show picture if it exist */}
                                            {ingInRecipe.ingredient.imageURL && (
                                                <img
                                                    src={ingInRecipe.ingredient.imageURL}
                                                    alt={ingInRecipe.ingredient.ingredient}
                                                    style={{
                                                        width: '40px',
                                                        height: '40px',
                                                        objectFit: 'cover',
                                                        marginRight: '15px',
                                                        borderRadius: '50%'
                                                    }}
                                                />
                                            )}
                                            <span>
                                            {ingInRecipe.quantity} {ingInRecipe.ingredient.unit} {ingInRecipe.ingredient.ingredient}
                                        </span>
                                        </li>
                                    ))}
                                </ul>
                            </div>
                            <div className="col-md-7">
                                <h3>Instructions</h3>
                                <ol>
                                    {recipe.instructions.sort((a, b) => a.stepSequence - b.stepSequence).map(inst => (
                                        <li key={inst.id} className="mb-2">{inst.stepDescription}</li>
                                    ))}
                                </ol>
                            </div>
                        </div>
                        {recipe.nutrition && (
                            <div className="mt-5">
                                <h3>Nutrition Facts</h3>
                                <table className="table table-striped">
                                    <tbody>
                                    <tr>
                                        <td>Calories</td>
                                        <td>{recipe.nutrition.calories} kcal</td>
                                    </tr>
                                    <tr>
                                        <td>Fat</td>
                                        <td>{recipe.nutrition.fat} g</td>
                                    </tr>
                                    <tr>
                                        <td>Carbohydrate</td>
                                        <td>{recipe.nutrition.carbohydrate} g</td>
                                    </tr>
                                    <tr>
                                        <td>Sugar</td>
                                        <td>{recipe.nutrition.sugar} g</td>
                                    </tr>
                                    <tr>
                                        <td>Fiber</td>
                                        <td>{recipe.nutrition.fiber} g</td>
                                    </tr>
                                    <tr>
                                        <td>Protein</td>
                                        <td>{recipe.nutrition.protein} g</td>
                                    </tr>
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                </div>
            </div>
            );
            }
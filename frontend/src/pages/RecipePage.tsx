import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

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
    };
    imageURL: string | null;
    portion: number;
    preparationTime: number | null;
    cookingTime: number | null;
    ingredients: IngredientInRecipe[];
    instructions: Instruction[];
    nutrition: Nutrition | null;
}

export default function RecipePage() {
    const { id } = useParams<{ id: string }>(); // resd from URL
    const [recipe, setRecipe] = useState<FullRecipe | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchRecipe = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/recipe/${id}`, {
                    credentials: 'include',
                });
                if (!response.ok) {
                    throw new Error(`Recipe not found or server error.`);
                }
                const data = await response.json();
                setRecipe(data);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchRecipe();
    }, [id]);

    if (loading) return <div className="container mt-4"><h2>Loading...</h2></div>;
    if (error) return <div className="container mt-4"><div className="alert alert-danger">{error}</div></div>;
    if (!recipe) return <div className="container mt-4"><h2>Recipe not found.</h2></div>;

    return (
        <div className="container mt-5">
            <div className="row">
                <div className="col-md-8 offset-md-2">
                    <h1>{recipe.title}</h1>
                    <p className="text-muted">By {recipe.author.name}</p>

                    {recipe.imageURL && <img src={recipe.imageURL} alt={recipe.title} className="img-fluid rounded mb-4" />}

                    <p className="lead">{recipe.description}</p>

                    <div className="my-4">
                        <span><strong>Servings:</strong> {recipe.portion} | </span>
                        <span><strong>Prep Time:</strong> {recipe.preparationTime === null ? '-' : `${recipe.preparationTime} min`} | </span>
                        <span><strong>Prep Time:</strong> {recipe.cookingTime === null ? '-' : `${recipe.cookingTime} min`} | </span>
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
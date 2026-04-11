import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';

interface Recipe {
    id: number;
    title: string;
    description: string;
    portion: number;
    preparationTime: number;
    cookingTime: number;
    published: boolean;
    imageURL: string;
}

export default function HomePage() {
    // to create link on the card
    const navigate = useNavigate();
    // State for saving the recipe list
    const [recipes, setRecipes] = useState<Recipe[]>([]);
    // State for loading
    const [loading, setLoading] = useState<boolean>(true);
    // State for possible errors
    const [error, setError] = useState<string | null>(null);

    // useEffect runs with first loading of the component.
    useEffect(() => {
        const fetchRecipes = async () => {
            try {
                const response = await api.get(`/recipes`);
                setRecipes(response.data);
            } catch (err) {
                setError('Failed to fetch recipes. Please try again later.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchRecipes();
    }, []);

    if (loading) {
        return <div className="container mt-4"><h2>Loading recipes...</h2></div>;
    }

    if (error) {
        return <div className="container mt-4"><div className="alert alert-danger">{error}</div></div>;
    }

    return (
        <div className="container mt-4">
            <h1>Find the best Recipes</h1>
            <div className="row">
                {recipes.length > 0 ? (
                    recipes.map(recipe => (
                        <div key={recipe.id} className="col-md-4 mb-4">
                            <div className="card h-100"
                                 onClick={() => navigate(`/recipe/${recipe.id}`)}
                                 style={{ cursor: 'pointer' }} >
                                {recipe.imageURL && (
                                    <img src={recipe.imageURL} className="recipe-card-img" alt={recipe.title} />
                                )}
                                <div className="card-body">
                                    <h5 className="card-title">{recipe.title}</h5>
                                    <p className="card-text">{recipe.description}</p>
                                    <p className="card-text">
                                        <small className="text-muted">
                                            Prep time: {recipe.preparationTime} min | Cook time: {recipe.cookingTime} min
                                        </small>
                                    </p>
                                </div>
                            </div>
                        </div>
                    ))
                ) : (
                    <p>No recipes found.</p>
                )}
            </div>
        </div>
    );
}
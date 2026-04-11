import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';

interface Recipe {
    id: number;
    title: string;
    description: string | null;
    portion: number;
    preparationTime: number | null;
    cookingTime: number | null;
    published: boolean;
    imageURL: string | null;
}

export default function MyRecipesPage() {
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
                const response = await api.get('/myrecipes');
                setRecipes(response.data);
            } catch (err:any) {
                //user is not logged in
                if (err.response && err.response.status === 401) {
                    setError('You must be logged in to view this page.');
                    setLoading(false);
                    return;
                } else if (err.response && err.response.status === 500) {
                    setError('Internal server error. Please try again later.');
                    console.error(err);
                    return;
                } else {
                    setError('Failed to connect to the server. Please try again later.');
                }
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
            <h1>My Recipes</h1>
            {/* User doesn't have recipes */}
            {(!recipes || recipes.length === 0) ? (
                <div className="text-center p-5 border rounded">
                    <h4>You haven't created any recipes yet.</h4>
                    <p>Why not create your first one?</p>
                    <a href="/createRecipe" className="btn btn-primary" role="button">Create recipe</a>
                </div>
            ) : (
                // Successfully showing the recipes
            <div className="row">
                {recipes.map(recipe => (
                    <div key={recipe.id} className="col-md-4 mb-4">
                        <div className="card h-100"
                             onClick={() => navigate(`/myrecipe/${recipe.id}`)}
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
                                    <p className="card-text">Recipe published: {recipe.published ? 'Yes' : 'No'} </p>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
            )}
        </div>
    );
}
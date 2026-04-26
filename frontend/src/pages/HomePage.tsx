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

    // states for filters and sorting
    const [sortOption, setSortOption] = useState('modified_at_desc');
    const [maxCalories, setMaxCalories] = useState('');
    const [minCalories, setMinCalories] = useState('');
    const [maxProteins, setMaxProteins] = useState('');
    const [minProteins, setMinProteins] = useState('');
    const [maxTotalTime, setMaxTotalTime] = useState('');
    // TODO: api to get list of authors


    // useEffect runs with first loading of the component.
    useEffect(() => {
        const fetchRecipes = async () => {
            setLoading(true);
            setError(null);

            try {

                //build URL parameters
                const params = new URLSearchParams();
                // add sorting
                params.append('sort', sortOption);
                // add filters if they have value
                if (maxCalories) {
                    params.append('max_calories', maxCalories);
                }
                if (minCalories) {
                    params.append('min_calories', minCalories);
                }
                if (maxProteins) {
                    params.append('max_proteins', maxProteins);
                }
                if (minProteins) {
                    params.append('min_proteins', minProteins);
                }
                if (maxTotalTime) {
                    params.append('max_total_time', maxTotalTime);
                }

                // const response = await api.get(`/recipes`);
                const response = await api.get(`/recipes?${params.toString()}`);
                setRecipes(response.data || []);
            } catch (err) {
                setError('Failed to fetch recipes. Please try again later.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchRecipes();
    }, [sortOption, maxCalories, maxTotalTime]);

    if (loading) {
        return <div className="container mt-4"><h2>Loading recipes...</h2></div>;
    }

    if (error) {
        return <div className="container mt-4"><div className="alert alert-danger">{error}</div></div>;
    }

    return (
        <div className="container mt-4">
            <h1>Find the best Recipes</h1>

            {/*filters and sorting*/}
            <div className="row g-3 align-items-center mb-4 p-3 bg-light rounded">
                <div className="col-md-4">
                    <label htmlFor="sort" className="form-label">Sort by</label>
                    <select
                        id="sort"
                        className="form-select"
                        value={sortOption}
                        onChange={(e) => setSortOption(e.target.value)}
                    >
                        <option value="modified_at_desc">Newest</option>
                        <option value="modified_at_asc">Oldest</option>
                        <option value="total_time_asc">Total Time (Asc)</option>
                        <option value="total_time_desc">Total Time (Desc)</option>
                        <option value="calories_asc">Calories (Asc)</option>
                        <option value="calories_desc">Calories (Desc)</option>
                        <option value="proteins_asc">Proteins (Asc)</option>
                        <option value="proteins_desc">Proteins (Desc)</option>

                    </select>
                </div>
                <div className="col-md-4">
                    <label htmlFor="maxCalories" className="form-label">Max Calories</label>
                    <input
                        type="number"
                        id="maxCalories"
                        className="form-control"
                        placeholder="e.g., 500"
                        value={maxCalories}
                        onChange={(e) => setMaxCalories(e.target.value)}
                    />
                </div>
                <div className="col-md-4">
                    <label htmlFor="minCalories" className="form-label">Min Calories</label>
                    <input
                        type="number"
                        id="minCalories"
                        className="form-control"
                        placeholder="e.g., 500"
                        value={minCalories}
                        onChange={(e) => setMinCalories(e.target.value)}
                    />
                </div>
                <div className="col-md-4">
                    <label htmlFor="maxProteins" className="form-label">Max Proteins</label>
                    <input
                        type="number"
                        id="maxProteins"
                        className="form-control"
                        placeholder="e.g., 500"
                        value={maxProteins}
                        onChange={(e) => setMaxProteins(e.target.value)}
                    />
                </div>
                <div className="col-md-4">
                    <label htmlFor="minProteins" className="form-label">Min Proteins</label>
                    <input
                        type="number"
                        id="minProteins"
                        className="form-control"
                        placeholder="e.g., 500"
                        value={minProteins}
                        onChange={(e) => setMinProteins(e.target.value)}
                    />
                </div>
                <div className="col-md-4">
                    <label htmlFor="maxTotalTime" className="form-label">Max Total Time (min)</label>
                    <input
                        type="number"
                        id="maxTotalTime"
                        className="form-control"
                        placeholder="e.g., 60"
                        value={maxTotalTime}
                        onChange={(e) => setMaxTotalTime(e.target.value)}
                    />
                </div>
            </div>

            {/*recipe cards*/}
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
// import { Link } from 'react-router-dom';
//
// export default function HomePage() {
//     return (
//         <div className="text-center">
//             <h1>Recipe Manager</h1>
//             <p className="lead">
//                 Save, organize and manage your recipes with ease.
//             </p>
//
//             <div className="mt-4">
//                 <b>What do you want too cook today?     </b>
//
//                 <input></input>
//                 <Link to="/recipes" className="btn btn-primary me-2">
//                     Find Recipes
//                 </Link>
//             </div>
//         </div>
//     );
// }


import { useState, useEffect } from 'react';

interface Recipe {
    id: number;
    title: string;
    description: string;
    portion: number;
    preparationTime: number;
    cookingTime: number;
    published: boolean;
}

export default function HomePage() {
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
                const response = await fetch('http://localhost:8080/api/recipes');
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: Recipe[] = await response.json();
                setRecipes(data);
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
                            <div className="card h-100">
                                <div className="card-body">
                                    <h5 className="card-title">{recipe.title}</h5>
                                    <p className="card-text">{recipe.description}</p>
                                    <p className="card-text">
                                        <small className="text-muted">
                                            Prep time: {recipe.preparationTime} min | Cook time: {recipe.cookingTime} min
                                        </small>
                                    </p>
                                    {/* placeholder to link to full recipe page */}
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
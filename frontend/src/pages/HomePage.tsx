import { Link } from 'react-router-dom';

export default function HomePage() {
    return (
        <div className="text-center">
            <h1>Recipe Manager</h1>
            <p className="lead">
                Save, organize and manage your recipes with ease.
            </p>

            <div className="mt-4">
                <b>What do you want too cook today?     </b>

                <input></input>
                <Link to="/recipes" className="btn btn-primary me-2">
                    Find Recipes
                </Link>
            </div>
        </div>
    );
}

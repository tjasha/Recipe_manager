export default function HomePage() {
    return (
        <div className="text-center">
            <h1>Recipe Manager</h1>
            <p className="lead">
                Save, organize and manage your recipes with ease.
            </p>

            <div className="mt-4">
                <a href="/recipes" className="btn btn-primary me-2">
                    Browse Recipes
                </a>
                <a href="/login" className="btn btn-outline-secondary">
                    Login
                </a>
            </div>
        </div>
    );
}

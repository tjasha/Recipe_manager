import { Link, useNavigate } from 'react-router-dom';
import { googleLogout } from '@react-oauth/google';
import "bootstrap/dist/js/bootstrap.bundle.min.js";

export default function NavigationBar() {
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            // First, notify the backend to destroy the session
            const res = await fetch('http://localhost:8080/api/auth/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    credentials: 'include',
                },
            });

            if (res.ok) {
                console.log('Backend logout successful.');
            } else {
                console.error('Backend logout failed.');
                // Even if backend logout fails, proceed with frontend logout
            }
        } catch (error) {
            console.error('Network error during logout:', error);
            // Proceed with frontend logout even if network fails
        }

        // Log out from Google on the frontend
        googleLogout();
        console.log('Logged out from Google on the frontend.');

        // Redirect to the login page
        navigate('/login');
    };

    return (
        <nav className="navbar navbar-expand-lg navbar-light bg-light bg-body-tertiary">
            <div className="container-fluid">
                <Link className="navbar-brand" to="/">Recipe Manager</Link>
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"></span>
                </button>

                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav ms-auto">
                        <li className="nav-item">
                            <Link className="nav-link" to="/">Recipes</Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to="/favourites">Favourites</Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to="/shoppingList">Shopping List</Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to="/myrecipes">My Recipes</Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to="/createRecipe">Create Recipe</Link>
                        </li>
                        <li className="nav-item">
                            <Link className="nav-link" to="/login">Login</Link>
                        </li>
                        <li className="nav-item">
                            <button className="nav-link btn btn-link" onClick={handleLogout}>Logout</button>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    );
}

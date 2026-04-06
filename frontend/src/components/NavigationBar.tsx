import { Link, useNavigate } from 'react-router-dom';
import { googleLogout } from '@react-oauth/google';
import "bootstrap/dist/js/bootstrap.bundle.min.js";

interface User {
    id: number;
    name: string;
    email: string;
    accessLevel: number;
}

interface NavigationBarProps {
    user: User | null;
    onLogout: () => void;
}

export default function NavigationBar({ user, onLogout }: { user: any | null, onLogout: () => void }) {
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
            if (!res.ok) {
                console.error('Backend logout failed.');
            }
        } catch (error) {
            console.error('Network error during logout:', error);
        } finally {
            // Proceed with frontend logout even if network fails
            googleLogout()
            onLogout();
            navigate('/login');

        }
    };

    return (

        <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
            <div className="container-fluid">
                <Link className="navbar-brand" to="/">RecipeManager</Link>
                <button
                    className="navbar-toggler"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target="#navbarNav"
                    aria-controls="navbarNav"
                    aria-expanded="false"
                    aria-label="Toggle navigation"
                >
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                        <li className="nav-item">
                            <Link className="nav-link" to="/">Home</Link>
                        </li>

                        {/* Tabs only for chefs */}
                        {user && (
                            <>
                                <li className="nav-item">
                                    <Link className="nav-link" to="/myrecipes">My Recipes</Link>
                                </li>
                                <li className="nav-item">
                                    <Link className="nav-link" to="/createRecipe">Create Recipe</Link>
                                </li>
                                {user.accessLevel === 0 && (
                                    <li className="nav-item">
                                        <Link className="nav-link" to="/admin">Admin Panel</Link>
                                    </li>
                                )}
                            </>
                        )}
                    </ul>

                    {/* Logout and welcome message */}
                    <ul className="navbar-nav">
                        {user ? (
                            // If user is logged in
                            <li className="nav-item dropdown">
                                <a
                                    className="nav-link dropdown-toggle"
                                    href="#"
                                    id="navbarDropdown"
                                    role="button"
                                    data-bs-toggle="dropdown"
                                    aria-expanded="false"
                                >
                                    Welcome, {user.name}
                                </a>
                                <ul className="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
                                    <li>
                                        <button className="dropdown-item" onClick={handleLogout}>
                                            Logout
                                        </button>
                                    </li>
                                </ul>
                            </li>
                        ) : (
                            // If user is not logged in
                            <li className="nav-item">
                                <Link className="nav-link" to="/login">Login</Link>
                            </li>
                        )}
                    </ul>
                </div>
            </div>
        </nav>
    );
}

import { Link } from 'react-router-dom';

export default function NavigationBar() {
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
                        <Link className="nav-link" to="/myRecipes">My Recipes</Link>
                    </li>
                    <li className="nav-item">
                        <Link className="nav-link" to="/login">Login</Link>
                    </li>
                    <li className="nav-item">
                        <Link className="nav-link" to="/logout">Logout</Link>
                    </li>
                </ul>
            </div>
        </div>
    </nav>
)
    ;
}

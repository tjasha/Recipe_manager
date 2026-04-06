import './App.css'
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";
import { useState } from 'react';
import HomePage from './pages/HomePage';
import MainLayout from "./layouts/MainLayout";
import LoginPage from "./pages/LoginPage";
import RecipePage from "./pages/RecipePage";
import NavigationBar from "./components/NavigationBar";
import MyRecipesPage from "./pages/MyRecipesPage";
import CreateRecipePage from "./pages/CreateRecipePage";
import MyRecipePage from "./pages/MyRecipePage";
import EditRecipePage from "./pages/EditRecipePage";

// User interface to follow if user is logged in
interface User {
    id: number;
    name: string;
    email: string;
    accessLevel: number;
}

function App() {

    // managing user state
    const [user, setUser] = useState<User | null>(null);
    const navigate = useNavigate();

    // function for log in page
    const handleLogin = (userData: User) => {
        setUser(userData);
        navigate('/myrecipes'); // navigate to user's recipes page
    };

    // function for navigation bar
    const handleLogout = () => {
        setUser(null); // Clean user's state (session is cleaned by backend)
    };

    return (
        <>
            {/* sending state and function for log out to the component */}
            <NavigationBar user={user} onLogout={handleLogout} />

            <Routes>
                <Route path="/" element={<HomePage></HomePage>} />
                <Route path="/recipe/:id" element={<RecipePage/>} />

                {/* Sending function for log in to the component */}
                <Route path="/login" element={<LoginPage onLoginSuccess={handleLogin} />} />
                {/*<Route path="/login" element={<LoginPage/>} />*/}

                {/* Protected paths - only shown if user is logged in */}
                {user && (
                    <>
                        <Route path="/myrecipes" element={<MyRecipesPage/>} />
                        <Route path="/myrecipe/:id" element={<MyRecipePage/>} />
                        <Route path="/createRecipe" element={<CreateRecipePage/>} />
                        <Route path="/myrecipe/:id/edit" element={<EditRecipePage/>} />
                        {/* <Route path="*" element={<NotFound />} />*/}
                    </>
                )}
            </Routes>
        </>
    )
}

export default App
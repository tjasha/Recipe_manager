import './App.css'
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";
import { useState } from 'react';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import HomePage from './pages/HomePage';
import LoginPage from "./pages/LoginPage";
import RecipePage from "./pages/RecipePage";
import NavigationBar from "./components/NavigationBar";
import MyRecipesPage from "./pages/MyRecipesPage";
import CreateRecipePage from "./pages/CreateRecipePage";
import MyRecipePage from "./pages/MyRecipePage";
import EditRecipePage from "./pages/EditRecipePage";
import RequireAuth from './components/RequireAuth';
import CheckSession from './components/CheckSession';
import ManageUsersPage from './pages/ManageUsersPage';
import RequireRole from './components/RequireRole';

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
    // new state to know is initial session check is done
    const [sessionChecked, setSessionChecked] = useState(false);
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

    // This function will be called by CheckSession
    const handleSessionChecked = () => {
        setSessionChecked(true);
    };

    // If the session hasn't been checked yet, we can show a loading spinner
    // or nothing. We also render CheckSession to trigger the verification.
    if (!sessionChecked) {
        return <CheckSession onSessionChecked={handleSessionChecked} onLogin={handleLogin} />;
    }

    return (
        <>
            <ToastContainer
                position="top-right"
                autoClose={3000}
                aria-label="Notifications"
                hideProgressBar={false}
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
            />

            {/* sending state and function for log out to the component */}
            <NavigationBar user={user} onLogout={handleLogout} />

            <Routes>
                {/*Public routes*/}
                <Route path="/" element={<HomePage></HomePage>} />
                <Route path="/recipe/:id" element={<RecipePage/>} />
                {/* Sending function for log in to the component */}
                <Route path="/login" element={<LoginPage onLoginSuccess={handleLogin} />} />

                {/* Protected paths - only shown if user is logged in */}
                {/*components are wrapped in RequireAuth guard*/}
                <Route
                    path="/myrecipes"
                    element={
                        <RequireAuth user={user}>
                            <MyRecipesPage />
                        </RequireAuth>
                    }
                />
                <Route
                    path="/myrecipe/:id"
                    element={
                        <RequireAuth user={user}>
                            <MyRecipePage />
                        </RequireAuth>
                    }
                />
                <Route
                    path="/createRecipe"
                    element={
                        <RequireAuth user={user}>
                            <CreateRecipePage />
                        </RequireAuth>
                    }
                />
                <Route
                    path="/myrecipe/:id/edit"
                    element={
                        <RequireAuth user={user}>
                            <EditRecipePage />
                        </RequireAuth>
                    }
                />

                {/* Admin paths - only shown if user is admin */}
                {/*components are wrapped in RequireRole guard*/}
                <Route
                    path="/admin/users"
                    element={
                        // we use both guards
                        <RequireAuth user={user}>
                            <RequireRole user={user} allowedRoles={[0]}>
                                <ManageUsersPage currentUser={user} />
                            </RequireRole>
                        </RequireAuth>
                    }
                />
            </Routes>
        </>
    );
}

export default App
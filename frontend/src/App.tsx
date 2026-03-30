import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from './pages/HomePage';
import MainLayout from "./layouts/MainLayout";
import LoginPage from "./pages/LoginPage";
import RecipePage from "./pages/RecipePage";
import NavigationBar from "./components/NavigationBar";
import MyRecipesPage from "./pages/MyRecipesPage";
import CreateRecipePage from "./pages/CreateRecipePage";
import MyRecipePage from "./pages/MyRecipePage";
import EditRecipePage from "./pages/EditRecipePage";

function App() {
    return (
        <>
        <NavigationBar />
            <Routes>
                <Route path="/" element={<HomePage></HomePage>} />
                <Route path="/recipe/:id" element={<RecipePage/>} />
                <Route path="/myrecipes" element={<MyRecipesPage/>} />
                <Route path="/myrecipe/:id" element={<MyRecipePage/>} />
                <Route path="/createRecipe" element={<CreateRecipePage/>} />
                <Route path="/myrecipe/:id/edit" element={<EditRecipePage/>} />
                <Route path="/login" element={<LoginPage/>} />
                {/*<Route path="*" element={<NotFound />} />*/}
            </Routes>
        </>
    )
}

export default App
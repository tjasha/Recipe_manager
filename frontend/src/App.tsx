import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from './pages/HomePage';
import MainLayout from "./layouts/MainLayout";
import LoginPage from "./pages/LoginPage";

function App() {
    return (
        <BrowserRouter>
            <MainLayout>
                <Routes>
                    <Route path="/" element={<HomePage></HomePage>} />
                    {/*<Route path="/recipes" element={<RecipesPage />} />*/}
                    <Route path="/login" element={<LoginPage/>} />
                    {/*<Route path="*" element={<NotFound />} />*/}
                </Routes>
            </MainLayout>
        </BrowserRouter>
    )
}

export default App
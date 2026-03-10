import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
// @ts-ignore
import HomePage from "./pages/HomePage.tsx";
// @ts-ignore
import LoginPage from "./pages/LoginPage.jsx";
// @ts-ignore
import MainLayout from "./layouts/MainLayout.tsx";


function App() {
    return (
        <BrowserRouter>
            <MainLayout>
                <Routes>
                    <Route path="/" element={<HomePage />} />
                    {/*<Route path="/recipes" element={<RecipesPage />} />*/}
                    <Route path="/login" element={<LoginPage />} />
                    {/*<Route path="*" element={<NotFound />} />*/}
                </Routes>
            </MainLayout>
        </BrowserRouter>
    )
}

export default App
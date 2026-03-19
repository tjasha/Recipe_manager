import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from './pages/HomePage';
import MainLayout from "./layouts/MainLayout";
import LoginPage from "./pages/LoginPage";
import RecipePage from "./pages/RecipePage";
import NavigationBar from "./components/NavigationBar";

function App() {
    return (
        <>
        <NavigationBar />
            <Routes>
                <Route path="/" element={<HomePage></HomePage>} />
                <Route path="/recipe/:id" element={<RecipePage/>} />
                <Route path="/login" element={<LoginPage/>} />
                {/*<Route path="*" element={<NotFound />} />*/}
            </Routes>
        </>
    )
}

export default App
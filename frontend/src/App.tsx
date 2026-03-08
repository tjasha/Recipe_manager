import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage.tsx";
import LoginPage from "./pages/LoginPage.tsx";
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

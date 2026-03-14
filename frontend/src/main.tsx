// import { StrictMode } from 'react'
// import { createRoot } from 'react-dom/client'
// import './index.css'
// import "bootstrap/dist/css/bootstrap.min.css";
// import App from './App.tsx'
//
// createRoot(document.getElementById('root')).render(
//   <StrictMode>
//       <App />
//   </StrictMode>,
// )
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import "bootstrap/dist/css/bootstrap.min.css";
import App from './App'
import { GoogleOAuthProvider } from '@react-oauth/google'

const CLIENT_ID = "544707640060-4olov7abecaemigsaqgciprslivq6bfk.apps.googleusercontent.com"

createRoot(document.getElementById('root')).render(
    <StrictMode>
        <GoogleOAuthProvider clientId={CLIENT_ID}>
            <App />
        </GoogleOAuthProvider>
    </StrictMode>,
)


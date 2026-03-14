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
import { GoogleOAuthProvider} from '@react-oauth/google'

const googleClientId = "544707640060-u21l0kiofbp0i1bdue59c8qodqugam46.apps.googleusercontent.com"

createRoot(document.getElementById('root')).render(
    <StrictMode>
        <GoogleOAuthProvider clientId={googleClientId}>
            <App />
        </GoogleOAuthProvider>
    </StrictMode>,
)


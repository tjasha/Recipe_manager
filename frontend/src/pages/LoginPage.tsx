import { GoogleLogin, googleLogout } from '@react-oauth/google';
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';

export default function LoginPage({ onLoginSuccess }: { onLoginSuccess: (user: any) => void }) {
    const [loginError, setLoginError] = useState<string | null>(null);
    const navigate = useNavigate();

    const handleGoogleLoginSuccess = async (credentialResponse: any) => {
        setLoginError(null); // Reset error state
        console.log("Received credential from Google:", credentialResponse);

        try {
            // Send the ID token to your backend for verification and session creation
            const res = await fetch('http://localhost:8080/api/auth/google/verify', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                // The backend expects an object with a "credential" field
                body: JSON.stringify({
                    credential: credentialResponse.credential,
                }),
            });

            if (res.ok) {
                const user = await res.json();
                onLoginSuccess(user);
                console.log('Backend login successful. User:', user);
                // The backend has set a session cookie.
            } else {
                // Handle errors from the backend (e.g., invalid token, server issue)
                console.error('Backend login failed.');
                setLoginError('Login failed on the server. Please try again.');
            }
        } catch (error) {
            // Handle network errors
            console.error('Network error during login:', error);
            setLoginError('A network error occurred. Please check your connection.');
        }
    };

    const handleGoogleLoginError = () => {
        console.log('Google Login Failed');
        setLoginError('Google login failed. Please try again.');
    };

    return (
        <div className="container" style={{ maxWidth: '400px', margin: '5rem auto' }}>
            <h1 className="text-center mb-4">Login</h1>

            {loginError && (
                <div className="alert alert-danger" role="alert">
                    {loginError}
                </div>
            )}

            <div className="d-flex justify-content-center mb-3">
                <GoogleLogin
                    size={"large"}
                    shape={"pill"}
                    onSuccess={handleGoogleLoginSuccess}
                    onError={handleGoogleLoginError}
                    auto_select={false} // this is preventing automatic sign-in
                />
            </div>
        </div>
    );
}

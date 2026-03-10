import { GoogleLogin } from '@react-oauth/google';

export default function LoginPage() {
    return (
        <div className="container" style={{ maxWidth: '400px', margin: '5rem auto' }}>
            <h1 className="text-center mb-4">Login</h1>

            <div className="d-flex justify-content-center mb-3">
                <GoogleLogin
                    onSuccess={(credentialResponse) => console.log(credentialResponse)}
                    onError={() => console.log("Login Failed")}
                />
            </div>

        </div>
    );
}

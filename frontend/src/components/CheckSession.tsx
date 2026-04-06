import { useEffect, useState } from 'react';

// This component will handle the initial session check
export default function CheckSession({ onSessionChecked, onLogin }: { onSessionChecked: () => void, onLogin: (user: any) => void }) {

    useEffect(() => {
        const verifySession = async () => {
            try {
                // We will create this new endpoint on the backend
                const response = await fetch('http://localhost:8080/api/auth/check-session', {
                    credentials: 'include',
                });

                if (response.ok) {
                    const userData = await response.json();
                    // If the session is valid, call the login handler from App.tsx
                    onLogin(userData);
                }
            } catch (error) {
                console.error("Session check failed:", error);
            } finally {
                // Signal to the App that the check is complete
                onSessionChecked();
            }
        };

        verifySession();
    }, [onLogin, onSessionChecked]);

    // This component doesn't render anything itself
    return null;
}
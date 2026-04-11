import { useEffect, useState } from 'react';
import api from '../api/client';

// This component will handle the initial session check
export default function CheckSession({ onSessionChecked, onLogin }: { onSessionChecked: () => void, onLogin: (user: any) => void }) {

    useEffect(() => {
        const verifySession = async () => {
            try {
                // We will create this new endpoint on the backend
                const response = await api.get('/auth/check-session');

                const userData = await response.data;
                // If the session is valid, call the login handler from App.tsx
                onLogin(userData);

            } catch (error) {
                // axios automatically throw errors for all statuses != 200
                console.log("Session check failed:", error);
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
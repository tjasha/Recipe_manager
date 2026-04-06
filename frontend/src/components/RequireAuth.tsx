import { useLocation, Navigate } from 'react-router-dom';

// This component will wrap our protected routes
export default function RequireAuth({ user, children }: { user: any | null, children: JSX.Element }) {
    const location = useLocation();

    if (!user) {
        // If the user is not logged in, redirect them to the login page.
        // We save the current location so we can send them back after they log in.
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    // If the user is logged in, render the child component (the actual page)
    return children;
}
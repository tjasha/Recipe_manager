import { Navigate } from 'react-router-dom';

interface User {
    accessLevel: number;
}

interface RequireRoleProps {
    user: User | null;
    allowedRoles: number[];
    children: JSX.Element;
}

export default function RequireRole({ user, allowedRoles, children }: RequireRoleProps) {
    // if user exist and the role is allowed, show the page
    if (user && allowedRoles.includes(user.accessLevel)) {
        return children; // show page
    }

    // if user exist but his role doesn't have permission
    if (user) {
        return <Navigate to="/" replace />; // redirect to home page
    }

    // if user isn't log in, redirect to login page
    return <Navigate to="/login" replace />;
}
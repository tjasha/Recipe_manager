import { useState, useEffect } from 'react';

// User interface from  API
interface User {
    id: number;
    name: string;
    email: string;
    accessLevel: number;
    state: string;
}

// recognise logged in user
interface CurrentUser {
    id: number;
}

export default function ManageUsersPage({ currentUser }: { currentUser: CurrentUser | null }) {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // --- Load users ---
    const fetchUsers = async () => {
        setLoading(true);
        try {
            const response = await fetch('http://localhost:8080/api/admin/users', {
                credentials: 'include',
            });

            if (response.status === 401 || response.status === 403) {
                throw new Error('You do not have permission to view this page.');
            }
            if (!response.ok) {
                throw new Error('Failed to fetch users.');
            }

            const data = await response.json();
            setUsers(data || []);
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // Load users on first render
    useEffect(() => {
        fetchUsers();
    }, []);

    // --- Delete user ---
    const handleDelete = async (userId: number, userName: string) => {
        if (!window.confirm(`Are you sure you want to delete user "${userName}"? This action cannot be undone.`)) {
            return;
        }

        try {
            const response = await fetch(`http://localhost:8080/api/admin/users/${userId}`, {
                method: 'DELETE',
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error('Failed to delete the user.');
            }

            // After successful deletion remove user from the state
            setUsers(prevUsers => prevUsers.filter(u => u.id !== userId));

        } catch (err: any) {
            alert(`Error: ${err.message}`);
        }
    };

    // --- Update user role or status  ---
    const handleUpdateUser = async (
        userId: number,
        field: 'access_level' | 'state',
        value: number | string) => {
        try {
            const payload = { [field]: value };
            const response = await fetch(`http://localhost:8080/api/admin/users/${userId}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify(payload),
            });

            if (!response.ok) {
                throw new Error('Failed to update user.');
            }

            // locally update state to reflect the change
            setUsers(prevUsers =>
                prevUsers.map(u =>{
                    if (u.id === userId) {
                        const updatedField = field === 'access_level' ? 'accessLevel' : 'state';
                        return { ...u, [updatedField]: value };
                    }
                    return u;
                })
            );
        } catch (err: any) {
            alert(`Error: ${err.message}`);
            // In case of error load data again to synchronise state
            fetchUsers();
        }
    };

    if (loading) {
        return <div className="container mt-4"><h2>Loading users...</h2></div>;
    }

    if (error) {
        return <div className="container mt-4"><div className="alert alert-danger">{error}</div></div>;
    }

    return (
        <div className="container mt-4">
            <h1>Manage Users</h1>
            <div className="table-responsive">
                <table className="table table-striped table-hover align-middle">
                    <thead className="table-dark">
                    <tr>
                        <th scope="col">ID</th>
                        <th scope="col">Name</th>
                        <th scope="col">Email</th>
                        <th scope="col">Role</th>
                        <th scope="col">Status</th>
                        <th scope="col">Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {users.map(user => (
                        <tr key={user.id}>
                            <th scope="row">{user.id}</th>
                            <td>{user.name}</td>
                            <td>{user.email}</td>
                            <td>
                                <select
                                    className="form-select form-select-sm"
                                    value={user.accessLevel}
                                    onChange={(e) => handleUpdateUser(user.id, 'access_level', Number(e.target.value))}
                                    disabled={!currentUser || user.id === currentUser.id}
                                    aria-label={`Role for ${user.name}`}
                                >
                                    <option value={0}>Admin</option>
                                    <option value={1}>Chef</option>
                                </select>
                            </td>
                            <td>
                                <select
                                    className="form-select form-select-sm"
                                    value={user.state}
                                    onChange={(e) => handleUpdateUser(user.id, 'state', e.target.value)}
                                    disabled={!currentUser || user.id === currentUser.id}
                                    aria-label={`Status for ${user.name}`}
                                >
                                    <option value={'active'}>Active</option>
                                    <option value={'deactivated'}>Deactivated</option>
                                </select>
                            </td>
                            <td>
                                {currentUser && user.id !== currentUser.id && (
                                    <button
                                        className="btn btn-sm btn-outline-danger"
                                        onClick={() => handleDelete(user.id, user.name)}
                                    >
                                        Delete
                                    </button>
                                )}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
import React, { useState, useEffect } from 'react';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import apiClient from '../api/apiClient';

export default function LoginPage() {
    const [username, setName] = useState('');
    const [password, setPassword] = useState('');
    const [companyId, setCompanyId] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);
    const navigate = useNavigate();
    const location = useLocation();

    useEffect(() => {
        const queryParams = new URLSearchParams(location.search);
        if (queryParams.get('signup') === 'success') {
            setSuccessMessage('Account created successfully! Please log in.');
        }
    }, [location]);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);
        setError(null);
        setSuccessMessage(null);

        const companyIdNum = parseInt(companyId);
        if (isNaN(companyIdNum)) {
            setError('Company ID must be a valid number');
            setIsSubmitting(false);
            return;
        }

        try {
            const response = await apiClient.post('login', {
                    username: username,
                    password: password,
                    company_id: companyIdNum
                });

            localStorage.setItem('authToken', response.token);
            navigate('/');

        } catch (err) {
            setError(err instanceof Error ? err.message : 'An unexpected error occurred.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-zinc-900 text-white">
            <div className="w-full max-w-md p-8 mx-auto bg-zinc-800 rounded-lg shadow-lg">
                <h1 className="text-3xl font-bold text-center mb-6">Login</h1>
                
                {error && <div className="bg-red-900 text-red-200 p-3 rounded-lg text-center mb-4">{error}</div>}
                {successMessage && <div className="bg-green-900 text-green-200 p-3 rounded-lg text-center mb-4">{successMessage}</div>}
                
                <form onSubmit={handleLogin}>
                    {/* Username Field */}
                    <div className="mb-4">
                        <label className="block text-zinc-400 mb-2" htmlFor="name">Username</label>
                        <input id="name" type="text" value={username} onChange={(e) => setName(e.target.value)} required className="w-full p-3 bg-zinc-700 rounded border border-zinc-600 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                    </div>
                    <div className="mb-6">
                        <label className="block text-zinc-400 mb-2" htmlFor="password">Password</label>
                        <input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required className="w-full p-3 bg-zinc-700 rounded border border-zinc-600 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                    </div>
                    {/* Company ID Field */}
                    <div className="mb-6">
                        <label className="block text-zinc-400 mb-2" htmlFor="companyId">Company ID</label>
                        <input id="companyId" type="text" value={companyId} onChange={(e) => setCompanyId(e.target.value)} required className="w-full p-3 bg-zinc-700 rounded border border-zinc-600 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                    </div>
                    <button type="submit" disabled={isSubmitting} className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded-lg focus:outline-none focus:shadow-outline transition duration-300 disabled:bg-blue-800">
                        {isSubmitting ? 'Logging in...' : 'Login'}
                    </button>
                </form>

                {/* Link to Sign Up Page */}
                <p className="text-center mt-6 text-zinc-400">
                    New employee? <Link to="/signup" className="text-blue-400 hover:underline">Create an account</Link>
                </p>
            </div>
        </div>
    );
}
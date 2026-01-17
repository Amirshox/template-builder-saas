'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { register } from '@/lib/api';

export default function RegisterPage() {
    const router = useRouter();
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const res = await register(email, password, name);
            localStorage.setItem('token', res.token);
            router.push('/');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to register');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex min-h-screen items-center justify-center bg-gray-50">
            <div className="w-full max-w-md space-y-8 rounded-lg bg-white p-6 shadow-md">
                <div className="text-center">
                    <h2 className="text-3xl font-bold tracking-tight text-gray-900">Create Account</h2>
                    <p className="mt-2 text-sm text-gray-600">
                        Or{' '}
                        <Link href="/login" className="font-medium text-blue-600 hover:text-blue-500">
                            sign in to existing account
                        </Link>
                    </p>
                </div>
                <form className="mt-8 space-y-6" onSubmit={handleRegister}>
                    {error && (
                        <div className="rounded-md bg-red-50 p-4 text-sm text-red-700">
                            {error}
                        </div>
                    )}
                    <div className="-space-y-px rounded-md shadow-sm">
                        <div>
                            <input
                                type="text"
                                required
                                className="relative block w-full rounded-t-md border-0 p-3 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
                                placeholder="Full Name"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                            />
                        </div>
                        <div>
                            <input
                                type="email"
                                required
                                className="relative block w-full border-0 p-3 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
                                placeholder="Email address"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                        <div>
                            <input
                                type="password"
                                required
                                className="relative block w-full rounded-b-md border-0 p-3 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
                                placeholder="Password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            disabled={loading}
                            className="group relative flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:opacity-50"
                        >
                            {loading ? 'Create Account' : 'Sign up'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}

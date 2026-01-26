'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { Mail, Lock, LogIn, ArrowLeft } from 'lucide-react';

export default function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const router = useRouter();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
            const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            if (!res.ok) {
                const errorData = await res.json().catch(() => ({ message: 'Invalid credentials' }));
                throw new Error(errorData.message || errorData.error || 'Invalid credentials');
            }

            const data = await res.json();
            localStorage.setItem('token', data.token);
            router.push('/dashboard');
        } catch (err: any) {
            setError(err.message || 'Login failed. Please try again.');
        } finally {
            setLoading(false);
        }
    };
    return (
        <div className="min-h-screen bg-background flex items-center justify-center p-4">
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-blue-500/5 blur-[120px] rounded-full -z-10"></div>

            <div className="w-full max-w-md">
                <div className="text-center mb-10">
                    <Link href="/" className="inline-flex items-center gap-2 mb-6 group">
                        <ArrowLeft className="text-slate-500 group-hover:text-blue-500 transition-colors" size={20} />
                        <Image src="/cortex_SaaS.png" alt="Cortex Logo" width={40} height={40} className="rounded" />
                        <span className="text-2xl font-bold tracking-tight text-white">Cortex</span>
                    </Link>
                    <h1 className="text-2xl font-bold text-white tracking-tight">Welcome back</h1>
                    <p className="text-slate-400 mt-2">Log in to your organization dashboard.</p>
                </div>

                <div className="bg-slate-800/40 border border-slate-700 p-8 rounded-2xl shadow-2xl backdrop-blur-xl">
                    {error && (
                        <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-xs font-bold uppercase tracking-widest text-center animate-in fade-in zoom-in duration-300">
                            {error}
                        </div>
                    )}
                    <form onSubmit={handleLogin} className="space-y-6">
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-300 ml-1">Email Address</label>
                            <div className="relative group">
                                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                <input
                                    type="email"
                                    placeholder="name@company.com"
                                    required
                                    className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </div>
                        </div>

                        <div className="space-y-2">
                            <div className="flex justify-between items-center ml-1">
                                <label className="text-sm font-bold text-slate-300">Password</label>
                                <Link href="/reset-password" title="reset-password-link" className="text-xs text-blue-500 hover:text-blue-400 font-medium">Forgot password?</Link>
                            </div>
                            <div className="relative group">
                                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                <input
                                    type="password"
                                    placeholder="••••••••"
                                    required
                                    className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </div>
                        </div>

                        <button
                            type="submit"
                            disabled={loading}
                            className="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-slate-800 text-white py-4 rounded-xl font-bold transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2"
                        >
                            {loading ? 'Authenticating...' : 'Log In'}
                            {!loading && <LogIn size={18} />}
                        </button>
                    </form>
                </div>

                <p className="text-center text-slate-500 mt-8 text-sm font-medium">
                    Don't have an account? <Link href="/register" className="text-blue-500 hover:text-blue-400 transition-colors">Sign up</Link>
                </p>
            </div>
        </div>
    );
}

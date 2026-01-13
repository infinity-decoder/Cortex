'use client';

import Link from 'next/link';
import Image from 'next/image';
import { Mail, Lock, LogIn, ArrowLeft } from 'lucide-react';

export default function Login() {
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
                    <form className="space-y-6">
                        <div className="space-y-2">
                            <label className="text-sm font-bold text-slate-300 ml-1">Email Address</label>
                            <div className="relative group">
                                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                <input
                                    type="email"
                                    placeholder="name@company.com"
                                    required
                                    className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
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
                                />
                            </div>
                        </div>

                        <button type="submit" className="w-full bg-blue-600 hover:bg-blue-500 text-white py-4 rounded-xl font-bold transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2">
                            Log In <LogIn size={18} />
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

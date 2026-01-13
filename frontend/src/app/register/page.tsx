'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { Mail, Lock, User, Building, ArrowRight, ShieldCheck } from 'lucide-react';

export default function Register() {
    const [step, setStep] = useState(1);
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        orgName: '',
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const router = useRouter();

    const handleNext = (e: React.FormEvent) => {
        e.preventDefault();
        setStep(2);
    };

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const res = await fetch('http://localhost:8080/api/v1/auth/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    email: formData.email,
                    password: formData.password,
                    org_name: formData.orgName
                })
            });

            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.error || 'Registration failed');
            }

            const data = await res.json();
            localStorage.setItem('token', data.token);
            router.push('/dashboard');
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-background flex items-center justify-center p-4">
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-blue-500/5 blur-[120px] rounded-full -z-10"></div>

            <div className="w-full max-w-md">
                <div className="text-center mb-10">
                    <Link href="/" className="inline-flex items-center gap-2 mb-6">
                        <Image src="/cortex_SaaS.png" alt="Cortex Logo" width={40} height={40} className="rounded" />
                        <span className="text-2xl font-bold tracking-tight text-white">Cortex</span>
                    </Link>
                    <h1 className="text-2xl font-bold text-white tracking-tight">Create your account</h1>
                    <p className="text-slate-400 mt-2">Secure access to your attack surface intel.</p>
                </div>

                <div className="bg-slate-800/40 border border-slate-700 p-8 rounded-2xl shadow-2xl backdrop-blur-xl">
                    <div className="flex gap-2 mb-8">
                        <div className={`h-1 flex-1 rounded-full ${step >= 1 ? 'bg-blue-500' : 'bg-slate-700'}`}></div>
                        <div className={`h-1 flex-1 rounded-full ${step >= 2 ? 'bg-blue-500' : 'bg-slate-700'}`}></div>
                    </div>

                    {error && (
                        <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-xs font-bold uppercase tracking-widest text-center animate-in fade-in zoom-in duration-300">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleNext} className="space-y-6">
                        {step === 1 ? (
                            <>
                                <div className="space-y-2">
                                    <label className="text-sm font-bold text-slate-300 ml-1">Full Name</label>
                                    <div className="relative group">
                                        <User className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                        <input
                                            type="text"
                                            placeholder="John Doe"
                                            required
                                            className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                                            value={formData.name}
                                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-bold text-slate-300 ml-1">Email Address</label>
                                    <div className="relative group">
                                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                        <input
                                            type="email"
                                            placeholder="johndoe@example.com"
                                            required
                                            className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                                            value={formData.email}
                                            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-bold text-slate-300 ml-1">Password</label>
                                    <div className="relative group">
                                        <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                        <input
                                            type="password"
                                            placeholder="••••••••"
                                            required
                                            className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                                            value={formData.password}
                                            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <button type="submit" className="w-full bg-blue-600 hover:bg-blue-500 text-white py-4 rounded-xl font-bold transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2">
                                    Continue <ArrowRight size={18} />
                                </button>
                            </>
                        ) : (
                            <>
                                <div className="space-y-2">
                                    <label className="text-sm font-bold text-slate-300 ml-1">Organization Name</label>
                                    <div className="relative group">
                                        <Building className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500 transition-colors" size={18} />
                                        <input
                                            type="text"
                                            placeholder="Acme Security Corp"
                                            required
                                            className="w-full bg-slate-900 border border-slate-700 rounded-xl py-3 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                                            value={formData.orgName}
                                            onChange={(e) => setFormData({ ...formData, orgName: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <div className="bg-blue-500/5 border border-blue-500/20 p-4 rounded-xl flex gap-3 items-start">
                                    <ShieldCheck className="text-blue-500 mt-0.5 flex-shrink-0" size={18} />
                                    <p className="text-xs text-slate-400 leading-relaxed">
                                        By creating an organization, you will be the primary **Owner** and responsible for all scan activities within this trust boundary.
                                    </p>
                                </div>

                                <div className="flex flex-col gap-3">
                                    <button
                                        type="button"
                                        onClick={handleRegister}
                                        disabled={loading}
                                        className="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-slate-800 text-white py-4 rounded-xl font-bold transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2"
                                    >
                                        {loading ? 'Creating Account...' : 'Start Monitoring'}
                                        {!loading && <ArrowRight size={18} />}
                                    </button>
                                    <button type="button" onClick={() => setStep(1)} disabled={loading} className="text-sm text-slate-500 hover:text-slate-300 transition-colors py-2 font-medium">
                                        Back to user details
                                    </button>
                                </div>
                            </>
                        )}
                    </form>
                </div>

                <p className="text-center text-slate-500 mt-8 text-sm font-medium">
                    Already have an account? <Link href="/login" className="text-blue-500 hover:text-blue-400 transition-colors">Log in</Link>
                </p>
            </div>
        </div>
    );
}

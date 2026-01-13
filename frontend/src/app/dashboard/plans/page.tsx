'use client';

import DashboardLayout from '@/components/DashboardLayout';
import { Check, Shield, Zap, Target, Star, CreditCard } from 'lucide-react';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

const tiers = [
    {
        name: "Discovery",
        price: "$0",
        current: true,
        features: ["Recursive DNS", "3 Domains", "Weekly Scans", "Email Alerts"],
        icon: Target,
        color: "text-slate-400"
    },
    {
        name: "Professional",
        price: "$490",
        period: "/mo",
        current: false,
        features: ["Unlimited Domains", "Continuous Scanning", "K8s/Docker Probing", "SLA Support", "API Discovery"],
        icon: Zap,
        color: "text-blue-500",
        highlight: true
    },
    {
        name: "Enterprise",
        price: "Custom",
        current: false,
        features: ["Multi-Org Support", "SAML/SSO", "Dedicated Infra", "Audit Logs", "SLA Guarantee"],
        icon: Star,
        color: "text-purple-500"
    }
];

export default function Plans() {
    const router = useRouter();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
        }
    }, []);

    return (
        <DashboardLayout>
            <div className="max-w-6xl mx-auto space-y-12">
                <div className="text-center md:text-left">
                    <h1 className="text-3xl font-extrabold text-white tracking-tight">Plans & Resource Allocation</h1>
                    <p className="text-slate-500 font-medium text-sm mt-1 uppercase tracking-widest leading-none">Scale your monitoring capacity</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    {tiers.map((tier) => (
                        <div key={tier.name} className={`relative bg-slate-900/40 border p-8 rounded-[2.5rem] flex flex-col transition-all overflow-hidden ${tier.highlight ? 'border-blue-500/50 ring-1 ring-blue-500/20' : 'border-slate-800'}`}>
                            {tier.highlight && (
                                <div className="absolute top-0 right-0 bg-blue-600 text-white px-4 py-1 text-[8px] font-black uppercase tracking-widest rounded-bl-xl shadow-xl">
                                    Recommended
                                </div>
                            )}

                            <div className="mb-8">
                                <tier.icon className={`${tier.color} mb-4`} size={32} />
                                <h3 className="text-xl font-bold text-white uppercase tracking-tight">{tier.name}</h3>
                                <div className="flex items-baseline gap-1 mt-2">
                                    <span className="text-4xl font-black text-white tracking-tighter">{tier.price}</span>
                                    {tier.period && <span className="text-sm font-bold text-slate-500">{tier.period}</span>}
                                </div>
                            </div>

                            <ul className="space-y-4 mb-10 flex-grow">
                                {tier.features.map((feat) => (
                                    <li key={feat} className="flex items-center gap-3">
                                        <div className="w-5 h-5 rounded-full bg-slate-800 flex items-center justify-center">
                                            <Check size={12} className="text-blue-500" />
                                        </div>
                                        <span className="text-xs font-bold text-slate-400">{feat}</span>
                                    </li>
                                ))}
                            </ul>

                            {tier.current ? (
                                <div className="w-full py-3 bg-slate-800 rounded-xl text-center text-[10px] font-black uppercase text-slate-500 tracking-widest cursor-default border border-slate-700">
                                    Current Active Plan
                                </div>
                            ) : (
                                <button className={`w-full py-3 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all active:scale-95 ${tier.highlight ? 'bg-blue-600 text-white hover:bg-blue-500 shadow-xl shadow-blue-500/20' : 'bg-slate-800 text-white hover:bg-slate-700'}`}>
                                    Upgrade Plan
                                </button>
                            )}
                        </div>
                    ))}
                </div>

                {/* Billing Summary */}
                <div className="bg-slate-900/40 border border-slate-800 rounded-3xl p-8 flex flex-col md:flex-row items-center justify-between gap-8">
                    <div className="flex items-center gap-6">
                        <div className="w-14 h-14 bg-slate-800 rounded-2xl flex items-center justify-center text-slate-400">
                            <CreditCard size={28} />
                        </div>
                        <div>
                            <h4 className="text-lg font-bold text-white mb-1 uppercase tracking-tight">Next Billing Cycle</h4>
                            <p className="text-sm text-slate-500 font-medium">Your next payment is scheduled for February 01, 2026.</p>
                        </div>
                    </div>
                    <button className="bg-slate-800 hover:bg-slate-700 text-white px-6 py-3 rounded-xl text-[10px] font-black uppercase tracking-widest border border-slate-700 transition-all">
                        Update Payment Method
                    </button>
                </div>
            </div>
        </DashboardLayout>
    );
}

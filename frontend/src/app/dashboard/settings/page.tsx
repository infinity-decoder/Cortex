'use client';

import { useState } from 'react';
import DashboardLayout from '@/components/DashboardLayout';
import {
    Building,
    Users,
    Globe,
    Shield,
    CreditCard,
    Bell,
    Plus,
    Trash2,
    Mail,
    MoreVertical,
    CheckCircle2
} from 'lucide-react';

const tabs = [
    { id: 'general', name: 'Organization', icon: Building },
    { id: 'team', name: 'Team & RBAC', icon: Users },
    { id: 'domains', name: 'Verified Domains', icon: Globe },
    { id: 'billing', name: 'Plans & Billing', icon: CreditCard },
];

const mockUsers = [
    { name: 'John Doe', email: 'john@cortex.sh', role: 'Owner', status: 'active' },
    { name: 'Alice Smith', email: 'alice@cortex.sh', role: 'Admin', status: 'active' },
    { name: 'Bob Wilson', email: 'bob@contractor.com', role: 'Analyst', status: 'pending' },
];

export default function Settings() {
    const [activeTab, setActiveTab] = useState('general');

    return (
        <DashboardLayout>
            <div className="max-w-5xl mx-auto space-y-10">
                <div>
                    <h1 className="text-3xl font-extrabold text-white tracking-tight">Organization Settings</h1>
                    <p className="text-slate-500 font-medium text-sm mt-1 uppercase tracking-widest leading-none">Manage team, scope, and billing preferences</p>
                </div>

                <div className="flex flex-col lg:flex-row gap-10">
                    {/* Navigation Sidebar */}
                    <aside className="w-full lg:w-64 space-y-1">
                        {tabs.map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-bold transition-all ${activeTab === tab.id
                                        ? 'bg-blue-600/10 text-blue-400 border border-blue-500/20 shadow-lg shadow-blue-500/5'
                                        : 'text-slate-500 hover:text-slate-300 hover:bg-slate-900/50'
                                    }`}
                            >
                                <tab.icon size={18} />
                                {tab.name}
                            </button>
                        ))}
                    </aside>

                    {/* Content Area */}
                    <div className="flex-1 space-y-12">
                        {activeTab === 'general' && (
                            <section className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-300">
                                <div className="bg-slate-900/40 border border-slate-800 p-8 rounded-3xl space-y-6">
                                    <h3 className="text-lg font-bold text-white uppercase tracking-tight">General Information</h3>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                        <div className="space-y-2">
                                            <label className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest ml-1">Org Name</label>
                                            <input className="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-sm font-medium text-white focus:outline-none focus:ring-1 focus:ring-blue-500/50" defaultValue="Cortex Security" />
                                        </div>
                                        <div className="space-y-2">
                                            <label className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest ml-1">Legal Entity ID (Optional)</label>
                                            <input className="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-sm font-medium text-white focus:outline-none focus:ring-1 focus:ring-blue-500/50" placeholder="CRUX-9901-X" />
                                        </div>
                                    </div>
                                    <button className="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-xl font-bold text-xs transition-all active:scale-95 shadow-lg shadow-blue-500/10">Save Changes</button>
                                </div>

                                <div className="bg-red-500/5 border border-red-500/20 p-8 rounded-3xl">
                                    <h3 className="text-lg font-bold text-red-500 uppercase tracking-tight mb-2">Danger Zone</h3>
                                    <p className="text-slate-500 text-xs mb-6 max-w-lg">
                                        Once you delete an organization, all verified domains, scan history, and delta alerts will be PERMANENTLY removed. This action cannot be undone.
                                    </p>
                                    <button className="bg-red-600/10 hover:bg-red-600 text-red-500 hover:text-white px-6 py-2.5 rounded-xl font-bold text-xs transition-all border border-red-600/20 active:scale-95">
                                        Delete Organization
                                    </button>
                                </div>
                            </section>
                        )}

                        {activeTab === 'team' && (
                            <section className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-300">
                                <div className="flex items-center justify-between">
                                    <h3 className="text-lg font-bold text-white uppercase tracking-tight">Global Team Management</h3>
                                    <button className="bg-slate-800 hover:bg-slate-700 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all border border-slate-700 flex items-center gap-2">
                                        <Plus size={14} /> Invite Member
                                    </button>
                                </div>
                                <div className="bg-slate-900/40 border border-slate-800 rounded-2xl overflow-hidden">
                                    <table className="w-full text-left">
                                        <thead>
                                            <tr className="border-b border-slate-800 text-[10px] uppercase font-extrabold text-slate-500 tracking-widest bg-slate-900/20">
                                                <th className="px-6 py-4">User</th>
                                                <th className="px-6 py-4">Role</th>
                                                <th className="px-6 py-4">Status</th>
                                                <th className="px-6 py-4 text-right">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-slate-800/50">
                                            {mockUsers.map((user, i) => (
                                                <tr key={i} className="hover:bg-slate-800/20 transition-colors group">
                                                    <td className="px-6 py-5">
                                                        <div className="flex items-center gap-3">
                                                            <div className="w-8 h-8 rounded-lg bg-blue-500/10 flex items-center justify-center text-blue-500 font-bold text-xs">
                                                                {user.name.charAt(0)}
                                                            </div>
                                                            <div>
                                                                <div className="text-sm font-bold text-white leading-none mb-1">{user.name}</div>
                                                                <div className="text-[10px] text-slate-500 font-medium">{user.email}</div>
                                                            </div>
                                                        </div>
                                                    </td>
                                                    <td className="px-6 py-5">
                                                        <span className="text-[10px] font-bold bg-slate-800 border border-slate-700 px-2 py-0.5 rounded text-slate-300 uppercase tracking-widest">{user.role}</span>
                                                    </td>
                                                    <td className="px-6 py-5">
                                                        {user.status === 'active' ? (
                                                            <div className="flex items-center gap-1.5 text-green-500 text-[10px] font-bold uppercase tracking-widest">
                                                                <CheckCircle2 size={12} /> Active
                                                            </div>
                                                        ) : (
                                                            <div className="flex items-center gap-1.5 text-yellow-500 text-[10px] font-bold uppercase tracking-widest">
                                                                <Mail size={12} /> Pending
                                                            </div>
                                                        )}
                                                    </td>
                                                    <td className="px-6 py-5 text-right">
                                                        <button className="text-slate-600 hover:text-white"><MoreVertical size={18} /></button>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </section>
                        )}

                        {activeTab === 'domains' && (
                            <section className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-300">
                                <div className="bg-slate-900/40 border border-slate-800 p-8 rounded-3xl space-y-6">
                                    <div className="flex items-center justify-between">
                                        <h3 className="text-lg font-bold text-white uppercase tracking-tight">Trust Boundary Management</h3>
                                        <button className="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-lg active:scale-95">Verify New Domain</button>
                                    </div>
                                    <div className="space-y-4">
                                        {[
                                            { domain: 'cortex.security', method: 'DNS TXT', date: 'Jan 01, 2026' },
                                            { domain: 'staging.cortex.sh', method: 'DNS TXT', date: 'Jan 10, 2026' }
                                        ].map((dom, i) => (
                                            <div key={i} className="flex items-center justify-between p-4 bg-slate-950 border border-slate-800 rounded-2xl group hover:border-slate-700 transition-all">
                                                <div className="flex items-center gap-4">
                                                    <div className="w-10 h-10 rounded-xl bg-green-500/10 flex items-center justify-center text-green-500">
                                                        <Globe size={18} />
                                                    </div>
                                                    <div>
                                                        <div className="text-sm font-bold text-white font-mono">{dom.domain}</div>
                                                        <div className="text-[10px] text-slate-500 font-medium">Verified via {dom.method} on {dom.date}</div>
                                                    </div>
                                                </div>
                                                <button className="text-slate-700 hover:text-red-500 transition-colors p-2"><Trash2 size={18} /></button>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </section>
                        )}
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
}

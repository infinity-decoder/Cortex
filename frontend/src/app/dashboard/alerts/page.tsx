'use client';

import DashboardLayout from '@/components/DashboardLayout';
import { Bell, AlertTriangle, ShieldCheck, Mail, MessageSquare, Plus, MoreVertical, Trash2 } from 'lucide-react';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

const mockAlerts = [
    { id: 'AL-001', type: 'Exposure Detected', target: 'cortex.security', severity: 'critical', desc: 'Kubernetes API server exposed on port 6443 with unauthenticated access.', date: '2h ago', status: 'unread' },
    { id: 'AL-002', type: 'Asset Delta', target: 'staging.cortex.sh', severity: 'medium', desc: 'New subdomain detected: api-dev.staging.cortex.sh (IP: 1.2.3.4)', date: '5h ago', status: 'read' },
    { id: 'AL-003', type: 'Service Change', target: 'cortex.security', severity: 'low', desc: 'Port 80/TCP closed on primary load balancer.', date: 'Yesterday', status: 'read' },
];

export default function Alerts() {
    const [alerts, setAlerts] = useState(mockAlerts);
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
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div>
                        <h1 className="text-3xl font-extrabold text-white tracking-tight">Security Alerts</h1>
                        <p className="text-slate-500 font-medium text-sm mt-1 uppercase tracking-widest leading-none">Real-time delta tracking & exposure notifications</p>
                    </div>
                    <div className="flex gap-4">
                        <button className="bg-slate-800 hover:bg-slate-700 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all border border-slate-700 flex items-center gap-2">
                            <Mail size={14} /> Configure SMTP
                        </button>
                        <button className="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-lg shadow-blue-500/10 flex items-center gap-2 active:scale-95">
                            <Plus size={14} /> New Integration
                        </button>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                    {/* Left: Summary Stats */}
                    <div className="lg:col-span-1 space-y-6">
                        <div className="bg-slate-900/40 border border-slate-800 p-6 rounded-2xl">
                            <h3 className="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-4 leading-none">Alert Summary</h3>
                            <div className="space-y-4">
                                <div className="flex justify-between items-center">
                                    <span className="text-sm font-bold text-slate-300">Critical</span>
                                    <span className="text-sm font-black text-red-500">1</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-sm font-bold text-slate-300">Medium</span>
                                    <span className="text-sm font-black text-orange-500">1</span>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-sm font-bold text-slate-300">Security Delta</span>
                                    <span className="text-sm font-black text-blue-500">12</span>
                                </div>
                            </div>
                        </div>

                        <div className="bg-blue-600/5 border border-blue-600/10 p-6 rounded-2xl">
                            <div className="flex items-center gap-3 mb-4">
                                <ShieldCheck className="text-blue-500" size={20} />
                                <h3 className="text-[10px] font-black text-blue-500 uppercase tracking-widest leading-none">Monitoring Status</h3>
                            </div>
                            <p className="text-xs text-slate-400 font-medium leading-relaxed">
                                Continuous verification is active for 4 verified domains.
                            </p>
                        </div>
                    </div>

                    {/* Right: Alert List */}
                    <div className="lg:col-span-3 space-y-4">
                        {alerts.map((alert) => (
                            <div key={alert.id} className={`group relative bg-slate-900/40 border p-6 rounded-2xl overflow-hidden transition-all hover:bg-slate-900/60 ${alert.status === 'unread' ? 'border-blue-500/30 ring-1 ring-blue-500/10' : 'border-slate-800'}`}>
                                <div className="flex items-start justify-between mb-4">
                                    <div className="flex items-center gap-4">
                                        <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${alert.severity === 'critical' ? 'bg-red-500/10 text-red-500' : 'bg-blue-500/10 text-blue-500'}`}>
                                            {alert.severity === 'critical' ? <AlertTriangle size={20} /> : <Bell size={20} />}
                                        </div>
                                        <div>
                                            <div className="flex items-center gap-3">
                                                <h4 className="text-base font-bold text-white tracking-tight">{alert.type}</h4>
                                                <span className={`text-[8px] font-black uppercase tracking-widest px-2 py-0.5 rounded-full border ${alert.severity === 'critical' ? 'bg-red-500/10 text-red-500 border-red-500/20' : 'bg-slate-800 text-slate-500 border-slate-700'}`}>
                                                    {alert.severity}
                                                </span>
                                            </div>
                                            <div className="text-[10px] text-slate-500 font-mono font-bold uppercase tracking-widest mt-1">
                                                Target: {alert.target} • {alert.id}
                                            </div>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <div className="text-[10px] font-bold text-slate-600 uppercase tracking-widest">{alert.date}</div>
                                        <button className="text-slate-700 hover:text-white transition-colors mt-2"><MoreVertical size={18} /></button>
                                    </div>
                                </div>
                                <p className="text-sm text-slate-400 font-medium leading-relaxed max-w-2xl ml-14">
                                    {alert.desc}
                                </p>
                                <div className="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <button className="text-slate-700 hover:text-red-500"><Trash2 size={16} /></button>
                                </div>
                            </div>
                        ))}

                        {alerts.length === 0 && (
                            <div className="py-20 text-center bg-slate-900/20 border border-slate-800 border-dashed rounded-3xl">
                                <ShieldCheck size={48} className="mx-auto text-slate-800 mb-4 opacity-20" />
                                <p className="text-sm text-slate-600 font-bold uppercase">No active alerts. Your surface is stable.</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
}

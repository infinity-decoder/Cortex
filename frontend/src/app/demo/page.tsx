'use client';

import { useState } from 'react';
import DashboardLayout from '@/components/DashboardLayout';
import {
    ShieldAlert,
    Globe,
    Zap,
    Activity,
    ArrowUpRight,
    Plus,
    Play,
    Info
} from 'lucide-react';

const mockAssets = [
    { Subdomain: 'api.cortex.security', IpAddress: '104.21.31.90' },
    { Subdomain: 'staging-k8s.cortex.security', IpAddress: '35.232.10.44' },
    { Subdomain: 'vpn.cortex.security', IpAddress: '142.250.190.46' },
    { Subdomain: 'db-admin.cortex.security', IpAddress: '34.120.55.12' },
];

const mockFindings = [
    { type: 'Unauthenticated Docker API', severity: 'critical', Detected: '2h ago' },
    { type: 'Exposed Kubelet (Read-Only)', severity: 'high', Detected: '5h ago' },
    { type: 'Expired SSL Certificate', severity: 'medium', Detected: 'Yesterday' },
];

export default function DemoPage() {
    const [domain, setDomain] = useState('');

    return (
        <DashboardLayout>
            <div className="max-w-7xl mx-auto space-y-10">

                {/* Demo Mode Banner */}
                <div className="bg-blue-600/10 border border-blue-500/20 p-4 rounded-2xl flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <Info size={18} className="text-blue-500" />
                        <span className="text-sm font-bold text-blue-400 uppercase tracking-widest bg-blue-500/10 px-2 py-1 rounded">Demo Mode</span>
                        <p className="text-xs text-slate-400 font-medium">Viewing simulated environment. Real data requires organization verification.</p>
                    </div>
                    <button className="bg-blue-600 text-white px-4 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-widest shadow-lg active:scale-95">Get Started</button>
                </div>

                {/* Header Section */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div>
                        <h1 className="text-3xl font-extrabold text-white tracking-tight mb-1">Security Overview</h1>
                        <p className="text-slate-500 text-sm font-medium uppercase tracking-widest">Global Attack Surface Status</p>
                    </div>
                    <div className="flex items-center gap-3">
                        <div className="bg-slate-900 border border-slate-800 rounded-xl p-1 flex gap-1">
                            <button className="px-4 py-2 bg-slate-800 text-white rounded-lg text-xs font-bold shadow-sm">24H</button>
                            <button className="px-4 py-2 text-slate-500 hover:text-slate-300 rounded-lg text-xs font-bold transition-colors">7D</button>
                            <button className="px-4 py-2 text-slate-500 hover:text-slate-300 rounded-lg text-xs font-bold transition-colors">30D</button>
                        </div>
                    </div>
                </div>

                {/* Scorecards */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl relative overflow-hidden group">
                        <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                            <Globe size={64} />
                        </div>
                        <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">Total Assets</div>
                        <div className="text-3xl font-extrabold text-white mb-2">124</div>
                        <div className="flex items-center gap-1.5 text-green-500 text-xs font-bold">
                            <ArrowUpRight size={14} />
                            <span>+2 New</span>
                        </div>
                    </div>

                    <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl relative overflow-hidden group">
                        <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                            <Zap size={64} />
                        </div>
                        <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">Exposed Services</div>
                        <div className="text-3xl font-extrabold text-white mb-2">12</div>
                        <div className="flex items-center gap-1.5 text-orange-500 text-xs font-bold">
                            <Activity size={14} className="animate-pulse" />
                            <span>4 Active Endpoints</span>
                        </div>
                    </div>

                    <div className="bg-red-500/5 border border-red-500/20 p-6 rounded-2xl relative overflow-hidden group">
                        <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                            <ShieldAlert size={64} className="text-red-500" />
                        </div>
                        <div className="text-[10px] font-extrabold text-red-500/50 uppercase tracking-widest mb-1">Critical Risks</div>
                        <div className="text-3xl font-extrabold text-red-500 mb-2">3</div>
                        <div className="flex items-center gap-1.5 text-red-400 text-xs font-bold">
                            <span>Attention Required</span>
                        </div>
                    </div>

                    <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl relative overflow-hidden group">
                        <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">System Risk Score</div>
                        <div className="text-3xl font-extrabold text-white mb-2">84<span className="text-slate-600 text-xl">/100</span></div>
                        <div className="flex items-center gap-1.5 text-slate-500 text-xs font-bold">
                            <span>Industry Avg: 65</span>
                        </div>
                    </div>
                </div>

                {/* Main Sections */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    <div className="lg:col-span-2 space-y-8">
                        <section className="bg-slate-900/40 border border-slate-800 rounded-2xl p-8 relative overflow-hidden">
                            <div className="relative z-10">
                                <h2 className="text-xl font-bold text-white mb-4 flex items-center gap-3">
                                    <Play size={20} className="text-blue-500 fill-current" />
                                    Trigger Manual Recon
                                </h2>
                                <div className="flex gap-4">
                                    <input
                                        type="text"
                                        placeholder="Enter domain (e.g., prod-api.cortex.security)"
                                        className="flex-1 bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-sm font-medium text-white placeholder:text-slate-600 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-all"
                                        value={domain}
                                        onChange={(e) => setDomain(e.target.value)}
                                    />
                                    <button
                                        className="bg-blue-600 hover:bg-blue-500 text-white font-bold py-3 px-8 rounded-xl transition-all shadow-lg active:scale-95 flex items-center gap-2"
                                    >
                                        Execute Scan
                                    </button>
                                </div>
                            </div>
                        </section>

                        <section>
                            <div className="flex items-center justify-between mb-4">
                                <h3 className="text-lg font-bold text-white flex items-center gap-3">
                                    <Globe size={20} className="text-slate-400" />
                                    Active Monitoring Subdomains
                                </h3>
                            </div>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                {mockAssets.map((asset, i) => (
                                    <div key={i} className="bg-slate-900/50 border border-slate-800 p-4 rounded-xl flex items-center justify-between group hover:border-slate-700 transition-colors">
                                        <div className="flex items-center gap-3">
                                            <div className="w-2 h-2 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]"></div>
                                            <div>
                                                <div className="text-sm font-bold text-white font-mono">{asset.Subdomain}</div>
                                                <div className="text-[10px] text-slate-500 font-mono mt-0.5">{asset.IpAddress}</div>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </section>
                    </div>

                    <div className="space-y-8">
                        <section className="bg-slate-900/40 border border-slate-800 rounded-2xl p-6">
                            <h3 className="text-sm font-extrabold text-slate-500 uppercase tracking-widest mb-6">Recent Delta Alerts</h3>
                            <div className="space-y-6">
                                {mockFindings.map((finding, i) => (
                                    <div key={i} className="flex gap-4 relative">
                                        <div className="flex flex-col items-center">
                                            <div className={`w-3 h-3 rounded-full ${finding.severity === 'critical' ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.5)]' :
                                                finding.severity === 'high' ? 'bg-orange-500' : 'bg-blue-500'
                                                }`}></div>
                                            {i < mockFindings.length - 1 && <div className="w-px flex-1 bg-slate-800 my-2"></div>}
                                        </div>
                                        <div className="-mt-1 pb-4">
                                            <div className="text-sm font-bold text-white leading-tight">{finding.type}</div>
                                            <div className="text-[10px] text-slate-600 font-medium mt-1 uppercase tracking-wider">{finding.Detected}</div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </section>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
}

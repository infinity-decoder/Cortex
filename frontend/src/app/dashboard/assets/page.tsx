'use client';

import { useState } from 'react';
import DashboardLayout from '@/components/DashboardLayout';
import {
    Globe,
    Server,
    Hash,
    Filter,
    Search,
    ArrowUpRight,
    MoreVertical,
    ShieldCheck,
    ShieldAlert,
    Zap
} from 'lucide-react';

const tabs = [
    { id: 'domains', name: 'Domains', icon: Globe },
    { id: 'ips', name: 'IP Addresses', icon: Hash },
    { id: 'services', name: 'Services', icon: Server },
];

const mockAssets = {
    domains: [
        { name: 'api.cortex.security', ip: '34.120.54.21', status: 'verified', tech: ['Go', 'NGINX'], lastSeen: '12m ago' },
        { name: 'k8s-console.cortex.security', ip: '35.190.21.44', status: 'verified', tech: ['Kubernetes', 'Traefik'], lastSeen: '1h ago' },
        { name: 'dev-db.cortex.security', ip: '10.0.0.45', status: 'unverified', tech: ['PostgreSQL'], lastSeen: '3h ago' },
        { name: 'staging-app.cortex.security', ip: '34.150.99.12', status: 'verified', tech: ['React', 'Node.js'], lastSeen: '5m ago' },
    ],
    services: [
        { port: 443, proto: 'HTTPS', asset: 'api.cortex.sh', service: 'NGINX 1.25.1', risk: 'low' },
        { port: 10250, proto: 'TCP', asset: 'k8s-node-01', service: 'Kubelet API', risk: 'critical' },
        { port: 2375, proto: 'TCP', asset: 'docker-host', service: 'Docker Remote API', risk: 'critical' },
        { port: 5432, proto: 'TCP', asset: 'db-master', service: 'PostgreSQL', risk: 'medium' },
    ]
};

export default function Assets() {
    const [activeTab, setActiveTab] = useState('domains');

    return (
        <DashboardLayout>
            <div className="max-w-7xl mx-auto space-y-8">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div>
                        <h1 className="text-3xl font-extrabold text-white tracking-tight">Asset Inventory</h1>
                        <p className="text-slate-500 font-medium text-sm mt-1 uppercase tracking-widest leading-none">Management of discovered infrastructure</p>
                    </div>
                    <div className="flex gap-2">
                        <button className="bg-slate-800 hover:bg-slate-700 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all border border-slate-700 flex items-center gap-2">
                            <Filter size={14} /> Filter
                        </button>
                        <button className="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-lg shadow-blue-500/10 flex items-center gap-2 active:scale-95">
                            <Zap size={14} className="fill-current" /> Trigger Discovery
                        </button>
                    </div>
                </div>

                {/* Custom Tabs */}
                <div className="bg-slate-900/50 p-1 border border-slate-800 rounded-2xl inline-flex gap-1">
                    {tabs.map((tab) => (
                        <button
                            key={tab.id}
                            onClick={() => setActiveTab(tab.id)}
                            className={`flex items-center gap-2 px-6 py-3 rounded-xl text-xs font-bold uppercase tracking-wider transition-all ${activeTab === tab.id
                                    ? 'bg-slate-800 text-blue-400 shadow-sm'
                                    : 'text-slate-500 hover:text-slate-300'
                                }`}
                        >
                            <tab.icon size={16} />
                            {tab.name}
                        </button>
                    ))}
                </div>

                {/* Table View */}
                <div className="bg-slate-900/40 border border-slate-800 rounded-2xl overflow-hidden shadow-2xl">
                    <div className="p-4 border-b border-slate-800 flex items-center justify-between">
                        <div className="relative">
                            <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
                            <input
                                type="text"
                                placeholder={`Search ${activeTab}...`}
                                className="bg-slate-950 border border-slate-800 rounded-lg py-1.5 pl-9 pr-4 text-xs focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-all w-64"
                            />
                        </div>
                        <span className="text-[10px] font-bold text-slate-600 uppercase">Showing {mockAssets.domains.length} Results</span>
                    </div>

                    <table className="w-full text-left">
                        <thead>
                            <tr className="border-b border-slate-800 text-[10px] uppercase font-extrabold text-slate-500 tracking-widest bg-slate-900/20">
                                <th className="px-6 py-4">Asset Name</th>
                                <th className="px-6 py-4">Context</th>
                                <th className="px-6 py-4">Technology Stack</th>
                                <th className="px-6 py-4">Verification</th>
                                <th className="px-6 py-4 text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-800/50">
                            {activeTab === 'domains' && mockAssets.domains.map((asset, i) => (
                                <tr key={i} className="hover:bg-slate-800/30 transition-colors group cursor-pointer">
                                    <td className="px-6 py-5">
                                        <div className="flex items-center gap-3">
                                            <div className="w-8 h-8 rounded-lg bg-blue-500/10 flex items-center justify-center text-blue-500">
                                                <Globe size={16} />
                                            </div>
                                            <div>
                                                <div className="text-sm font-bold text-white font-mono leading-none mb-1 group-hover:text-blue-400 transition-colors">{asset.name}</div>
                                                <div className="text-[10px] text-slate-500 font-mono font-medium">{asset.ip}</div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="text-xs text-slate-400 font-medium">Last seen {asset.lastSeen}</div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="flex gap-2">
                                            {asset.tech.map((t, j) => (
                                                <span key={j} className="text-[10px] font-bold bg-slate-800 border border-slate-700 px-2 py-0.5 rounded text-slate-300">
                                                    {t}
                                                </span>
                                            ))}
                                        </div>
                                    </td>
                                    <td className="px-6 py-5">
                                        {asset.status === 'verified' ? (
                                            <div className="inline-flex items-center gap-1.5 text-[10px] font-bold text-green-500 bg-green-500/5 border border-green-500/20 px-2 py-1 rounded-full uppercase tracking-widest">
                                                <ShieldCheck size={12} /> Verified
                                            </div>
                                        ) : (
                                            <div className="inline-flex items-center gap-1.5 text-[10px] font-bold text-slate-500 bg-slate-500/5 border border-slate-500/20 px-2 py-1 rounded-full uppercase tracking-widest">
                                                Unverified
                                            </div>
                                        )}
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="flex justify-end items-center gap-3">
                                            <button className="text-slate-500 hover:text-white transition-colors">
                                                <ArrowUpRight size={18} />
                                            </button>
                                            <button className="text-slate-500 hover:text-white transition-colors">
                                                <MoreVertical size={18} />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                            {activeTab === 'services' && mockAssets.services.map((service, i) => (
                                <tr key={i} className="hover:bg-slate-800/30 transition-colors group cursor-pointer">
                                    <td className="px-6 py-5">
                                        <div className="flex items-center gap-3">
                                            <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${service.risk === 'critical' ? 'bg-red-500/10 text-red-500' :
                                                    service.risk === 'medium' ? 'bg-orange-500/10 text-orange-500' :
                                                        'bg-blue-500/10 text-blue-500'
                                                }`}>
                                                {service.risk === 'critical' ? <ShieldAlert size={16} /> : <Server size={16} />}
                                            </div>
                                            <div>
                                                <div className="text-sm font-bold text-white font-mono leading-none mb-1 group-hover:text-blue-400 transition-colors">{service.asset}</div>
                                                <div className="text-[10px] text-slate-500 font-mono font-medium">Port {service.port} ({service.proto})</div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="text-xs text-slate-400 font-bold uppercase tracking-tight">{service.service}</div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <span className={`text-[10px] font-extrabold uppercase tracking-widest px-2 py-1 rounded-full border ${service.risk === 'critical' ? 'bg-red-500/10 text-red-500 border-red-500/20' :
                                                service.risk === 'medium' ? 'bg-orange-500/10 text-orange-500 border-orange-500/20' :
                                                    'bg-blue-500/10 text-blue-500 border-blue-500/20'
                                            }`}>
                                            {service.risk} Risk
                                        </span>
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="text-[10px] font-bold text-slate-600">PRODUCTION</div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="flex justify-end items-center gap-3">
                                            <button className="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-[10px] text-white font-bold rounded-lg border border-slate-700 transition-all active:scale-95">Analyze</button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>

                    {activeTab === 'ips' && (
                        <div className="p-20 text-center">
                            <Hash size={48} className="mx-auto text-slate-800 mb-4" />
                            <h3 className="text-lg font-bold text-slate-500">IP Scanning Ready</h3>
                            <p className="text-sm text-slate-700 max-w-xs mx-auto mt-2">Associate unmapped IP space to discover hidden infrastructure components.</p>
                            <button className="mt-6 text-xs font-bold text-blue-500 hover:text-blue-400 uppercase tracking-widest transition-colors">Add IP Range</button>
                        </div>
                    )}
                </div>
            </div>
        </DashboardLayout>
    );
}

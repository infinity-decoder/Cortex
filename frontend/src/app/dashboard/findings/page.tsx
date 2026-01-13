'use client';

import { useState } from 'react';
import DashboardLayout from '@/components/DashboardLayout';
import {
    ShieldAlert,
    ShieldCheck,
    Info,
    ExternalLink,
    ChevronDown,
    Search,
    Filter,
    CheckCircle2,
    AlertOctagon,
    Copy,
    Terminal
} from 'lucide-react';

const mockFindings = [
    {
        id: 'F-102',
        title: 'Exposed Docker Remote API (Unauthenticated)',
        severity: 'critical',
        category: 'RCE / System Takeover',
        asset: 'docker.cortex.security',
        service: 'Docker API (2375/tcp)',
        detected: '4h ago',
        status: 'open',
        description: 'The Docker Remote API is accessible without authentication. Attackers can execute commands, pull/push images, and gain root access to the host system.',
        remediation: 'Disable TCP access to the Docker API or enforce MTLS authentication using certificates.',
        evidence: 'GET /v1.24/version HTTP/1.1\nHost: docker.cortex.security:2375\n\nHTTP/1.1 200 OK\nContent-Type: application/json\nApi-Version: 1.24\nArch: amd64\nOs: linux\n...'
    },
    {
        id: 'F-105',
        title: 'Kubernetes Kubelet API Anonymous Access',
        severity: 'critical',
        category: 'Privilege Escalation',
        asset: 'k8s-node-02',
        service: 'Kubelet (10250/tcp)',
        detected: '12h ago',
        status: 'open',
        description: 'The Kubelet API allows anonymous users to list pods and potentially execute commands via the /exec endpoint depending on specific cluster permissions.',
        remediation: 'Set --anonymous-auth=false and --authorization-mode=Webhook in Kubelet configuration.',
        evidence: 'curl -k https://node-ip:10250/pods\n\n{\n  "kind": "PodList",\n  "items": [...] \n}'
    },
    {
        id: 'F-201',
        title: 'TLS Certificate Set to Expire Soon',
        severity: 'medium',
        category: 'Service Availability',
        asset: 'api.cortex.security',
        service: 'HTTPS (443/tcp)',
        detected: '1d ago',
        status: 'acknowledged',
        description: 'The TLS certificate for this asset will expire in less than 7 days. This will cause service disruption and browser warnings.',
        remediation: 'Renew the certificate via Let\'s Encrypt or your primary CA provider.',
        evidence: 'Not Before: Jan 13 00:00:00 2026 GMT\nNot After : Jan 20 23:59:59 2026 GMT'
    }
];

export default function Findings() {
    const [selectedFinding, setSelectedFinding] = useState<any>(null);

    const getSeverityColor = (sev: string) => {
        switch (sev) {
            case 'critical': return 'text-red-500 bg-red-500/10 border-red-500/20';
            case 'high': return 'text-orange-500 bg-orange-500/10 border-orange-500/20';
            case 'medium': return 'text-yellow-500 bg-yellow-500/10 border-yellow-500/20';
            default: return 'text-blue-500 bg-blue-500/10 border-blue-500/20';
        }
    };

    return (
        <DashboardLayout>
            <div className="max-w-7xl mx-auto flex flex-col h-full">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                    <div>
                        <h1 className="text-3xl font-extrabold text-white tracking-tight">Security Findings</h1>
                        <p className="text-slate-500 font-medium text-sm mt-1 uppercase tracking-widest leading-none">Vulnerabilities and exposures detected</p>
                    </div>
                    <div className="flex gap-2">
                        <div className="bg-slate-900 border border-slate-800 rounded-xl p-1 flex gap-1">
                            <button className="px-4 py-2 bg-slate-800 text-white rounded-lg text-[10px] font-extrabold uppercase tracking-widest shadow-sm">All</button>
                            <button className="px-4 py-2 text-slate-500 hover:text-slate-300 rounded-lg text-[10px] font-extrabold uppercase tracking-widest transition-colors">Open</button>
                            <button className="px-4 py-2 text-slate-500 hover:text-slate-300 rounded-lg text-[10px] font-extrabold uppercase tracking-widest transition-colors">Fixed</button>
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 h-full min-h-[600px]">
                    {/* Findings List */}
                    <div className="space-y-4 overflow-y-auto pr-2 custom-scrollbar">
                        <div className="sticky top-0 bg-slate-950 z-10 pb-4 flex gap-4">
                            <div className="relative flex-1 group">
                                <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500 group-focus-within:text-blue-500" />
                                <input
                                    type="text"
                                    placeholder="Filter by ID, Title, or Category..."
                                    className="w-full bg-slate-900 border border-slate-800 rounded-xl py-2 pl-9 pr-4 text-xs focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-all font-medium"
                                />
                            </div>
                            <button className="bg-slate-900 border border-slate-800 p-2 rounded-xl text-slate-500">
                                <Filter size={18} />
                            </button>
                        </div>

                        {mockFindings.map((finding) => (
                            <div
                                key={finding.id}
                                onClick={() => setSelectedFinding(finding)}
                                className={`p-5 rounded-2xl border transition-all cursor-pointer group relative overflow-hidden ${selectedFinding?.id === finding.id
                                        ? 'bg-slate-800/80 border-blue-500/50 shadow-xl shadow-blue-500/5'
                                        : 'bg-slate-900/40 border-slate-800 hover:border-slate-700'
                                    }`}
                            >
                                <div className="flex justify-between items-start mb-3">
                                    <div className="flex items-center gap-2">
                                        <span className={`text-[10px] font-extrabold border px-2 py-0.5 rounded-full uppercase tracking-widest ${getSeverityColor(finding.severity)}`}>
                                            {finding.severity}
                                        </span>
                                        <span className="text-[10px] text-slate-500 font-bold font-mono tracking-tight">{finding.id}</span>
                                    </div>
                                    <span className="text-[10px] text-slate-600 font-bold uppercase tracking-widest">{finding.detected}</span>
                                </div>
                                <h3 className={`text-sm md:text-base font-bold text-white transition-colors group-hover:text-blue-400 ${selectedFinding?.id === finding.id ? 'text-blue-400' : ''}`}>
                                    {finding.title}
                                </h3>
                                <div className="flex items-center gap-4 mt-3">
                                    <div className="flex items-center gap-1.5">
                                        <Globe size={12} className="text-slate-600" />
                                        <span className="text-[10px] text-slate-500 font-mono">{finding.asset}</span>
                                    </div>
                                </div>
                                {selectedFinding?.id === finding.id && (
                                    <div className="absolute top-0 right-0 w-1 h-full bg-blue-500"></div>
                                )}
                            </div>
                        ))}
                    </div>

                    {/* Finding Detail Sidebar */}
                    <div className="bg-slate-900/40 border border-slate-800 rounded-3xl p-8 sticky top-0 h-fit">
                        {!selectedFinding ? (
                            <div className="h-full min-h-[400px] flex flex-col items-center justify-center text-center opacity-40">
                                <ShieldCheck size={64} className="mb-4 text-slate-800" />
                                <h3 className="text-lg font-bold text-slate-500">Select a finding</h3>
                                <p className="text-sm text-slate-700 mt-2 max-w-xs">Detailed exploitability evidence and remediation steps will appear here.</p>
                            </div>
                        ) : (
                            <div className="space-y-8 animate-in fade-in slide-in-from-right-4 duration-300">
                                <div className="flex justify-between items-start">
                                    <div>
                                        <span className="text-[10px] text-blue-500 font-extrabold uppercase tracking-[0.2em] mb-2 block">Exposure Intelligence</span>
                                        <h2 className="text-xl font-extrabold text-white leading-tight">{selectedFinding.title}</h2>
                                    </div>
                                    <button className="text-slate-500 hover:text-white"><CheckCircle2 size={24} /></button>
                                </div>

                                <div className="grid grid-cols-2 gap-4">
                                    <div className="bg-slate-950/50 p-4 rounded-2xl border border-slate-800">
                                        <div className="text-[10px] font-extrabold text-slate-600 uppercase mb-1">Attack Category</div>
                                        <div className="text-sm font-bold text-blue-400">{selectedFinding.category}</div>
                                    </div>
                                    <div className="bg-slate-950/50 p-4 rounded-2xl border border-slate-800">
                                        <div className="text-[10px] font-extrabold text-slate-600 uppercase mb-1">Status</div>
                                        <div className="text-sm font-bold text-green-500 flex items-center gap-1.5 uppercase tracking-widest">
                                            <div className="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
                                            {selectedFinding.status}
                                        </div>
                                    </div>
                                </div>

                                <div>
                                    <h4 className="text-xs font-bold text-slate-500 uppercase tracking-widest mb-3 flex items-center gap-2">
                                        <Info size={14} /> Description
                                    </h4>
                                    <p className="text-sm text-slate-400 leading-relaxed">
                                        {selectedFinding.description}
                                    </p>
                                </div>

                                <div>
                                    <h4 className="text-xs font-bold text-slate-500 uppercase tracking-widest mb-3 flex items-center gap-2">
                                        <ShieldCheck size={14} className="text-green-500" /> Remediation
                                    </h4>
                                    <div className="bg-slate-950 p-4 rounded-2xl border border-slate-800 border-l-4 border-l-green-500">
                                        <p className="text-sm text-slate-300 font-medium leading-relaxed italic">
                                            "{selectedFinding.remediation}"
                                        </p>
                                    </div>
                                </div>

                                <div>
                                    <div className="flex justify-between items-center mb-3">
                                        <h4 className="text-xs font-bold text-slate-500 uppercase tracking-widest flex items-center gap-2">
                                            <Terminal size={14} /> Technical Evidence
                                        </h4>
                                        <button className="text-slate-600 hover:text-blue-500 transition-colors"><Copy size={14} /></button>
                                    </div>
                                    <div className="bg-slate-950 border border-slate-800 rounded-2xl p-4 overflow-x-auto">
                                        <pre className="text-[11px] font-mono text-blue-300/80 leading-relaxed">
                                            {selectedFinding.evidence}
                                        </pre>
                                    </div>
                                </div>

                                <div className="pt-4 border-t border-slate-800 flex gap-4">
                                    <button className="flex-1 bg-blue-600 hover:bg-blue-500 text-white font-bold py-3 rounded-xl shadow-lg shadow-blue-500/10 active:scale-95 transition-all text-sm">
                                        Mark as Fixed
                                    </button>
                                    <button className="px-6 py-3 bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold rounded-xl border border-slate-700 transition-all text-sm">
                                        Ignore
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
}

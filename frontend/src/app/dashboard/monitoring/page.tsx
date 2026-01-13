'use client';

import { useState } from 'react';
import DashboardLayout from '@/components/DashboardLayout';
import {
    Clock,
    Calendar,
    RefreshCcw,
    Plus,
    Play,
    CheckCircle2,
    AlertCircle,
    MoreVertical,
    Activity,
    History,
    Zap
} from 'lucide-react';

const mockSchedules = [
    { id: 'SCH-01', target: 'cortex.security', frequency: 'Daily', lastRun: '2h ago', status: 'active', nextRun: '22h' },
    { id: 'SCH-02', target: 'staging.cortex.sh', frequency: 'Hourly', lastRun: '15m ago', status: 'active', nextRun: '45m' },
];

const mockHistory = [
    { id: 'RUN-998', target: 'cortex.security', type: 'Discovery', findings: 12, delta: '+2', status: 'complete', date: 'Jan 13, 08:30' },
    { id: 'RUN-997', target: 'staging.cortex.sh', type: 'Deep Scan', findings: 5, delta: '0', status: 'complete', date: 'Jan 13, 07:45' },
    { id: 'RUN-996', target: 'k8s-cluster-01', type: 'Audit', findings: 42, delta: '+8', status: 'error', date: 'Jan 13, 05:00' },
];

export default function Monitoring() {
    return (
        <DashboardLayout>
            <div className="max-w-7xl mx-auto space-y-12">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div>
                        <h1 className="text-3xl font-extrabold text-white tracking-tight">Monitoring & Schedules</h1>
                        <p className="text-slate-500 font-medium text-sm mt-1 uppercase tracking-widest leading-none">Continuous threat exposure management</p>
                    </div>
                    <button className="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-lg shadow-blue-500/10 flex items-center gap-2 active:scale-95">
                        <Plus size={14} /> Create Schedule
                    </button>
                </div>

                {/* Active Schedules */}
                <section>
                    <div className="flex items-center gap-3 mb-6">
                        <Calendar size={20} className="text-blue-500" />
                        <h2 className="text-lg font-bold text-white uppercase tracking-tight">Active Scan Schedules</h2>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {mockSchedules.map((sch) => (
                            <div key={sch.id} className="bg-slate-900/40 border border-slate-800 p-6 rounded-2xl group hover:border-blue-500/30 transition-all relative overflow-hidden">
                                <div className="flex justify-between items-start mb-6">
                                    <div className="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center text-blue-500">
                                        <Clock size={20} />
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <span className="text-[10px] font-bold text-green-500 bg-green-500/5 border border-green-500/20 px-2 py-0.5 rounded-full uppercase tracking-widest leading-normal">
                                            Active
                                        </span>
                                        <button className="text-slate-600 hover:text-white transition-colors"><MoreVertical size={18} /></button>
                                    </div>
                                </div>
                                <h3 className="text-base font-bold text-white mb-1 group-hover:text-blue-400 transition-colors">{sch.target}</h3>
                                <div className="text-[10px] text-slate-500 font-mono font-bold uppercase tracking-widest mb-4">{sch.id} • {sch.frequency}</div>

                                <div className="flex items-center justify-between pt-4 border-t border-slate-800/50">
                                    <div className="flex flex-col">
                                        <span className="text-[10px] font-bold text-slate-600 uppercase">Last Run</span>
                                        <span className="text-xs text-slate-300 font-medium">{sch.lastRun}</span>
                                    </div>
                                    <div className="flex flex-col text-right">
                                        <span className="text-[10px] font-bold text-slate-600 uppercase">Next Run In</span>
                                        <span className="text-xs text-blue-400 font-bold">{sch.nextRun}</span>
                                    </div>
                                </div>
                                {/* Decorative Pattern */}
                                <div className="absolute top-0 right-0 p-2 opacity-5 pointer-events-none">
                                    <RefreshCcw size={48} className="animate-spin-slow" />
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* Scan History & Delta Timeline */}
                <section>
                    <div className="flex items-center gap-3 mb-6">
                        <History size={20} className="text-slate-400" />
                        <h2 className="text-lg font-bold text-white uppercase tracking-tight">Scan Execution History</h2>
                    </div>
                    <div className="bg-slate-900/40 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
                        <table className="w-full text-left">
                            <thead>
                                <tr className="border-b border-slate-800 text-[10px] uppercase font-extrabold text-slate-500 tracking-widest bg-slate-900/20">
                                    <th className="px-6 py-4">Execution ID</th>
                                    <th className="px-6 py-4">Target Information</th>
                                    <th className="px-6 py-4">Scan Type</th>
                                    <th className="px-6 py-4">Findings (Delta)</th>
                                    <th className="px-6 py-4">Status</th>
                                    <th className="px-6 py-4 text-right">Evidence</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-800/50">
                                {mockHistory.map((run) => (
                                    <tr key={run.id} className="hover:bg-slate-800/20 transition-colors group cursor-pointer">
                                        <td className="px-6 py-5">
                                            <span className="text-xs font-bold text-slate-400 font-mono tracking-tight">{run.id}</span>
                                            <div className="text-[10px] text-slate-600 font-medium mt-1 uppercase">{run.date}</div>
                                        </td>
                                        <td className="px-6 py-5">
                                            <div className="text-sm font-bold text-white font-mono group-hover:text-blue-400 transition-colors">{run.target}</div>
                                        </td>
                                        <td className="px-6 py-5">
                                            <span className="text-[10px] font-bold bg-slate-800 border border-slate-700 px-2 py-0.5 rounded text-slate-400 uppercase tracking-widest">
                                                {run.type}
                                            </span>
                                        </td>
                                        <td className="px-6 py-5">
                                            <div className="flex items-center gap-2">
                                                <span className="text-sm font-extrabold text-white">{run.findings}</span>
                                                {run.delta !== '0' && (
                                                    <span className={`text-[10px] font-bold px-1.5 py-0.5 rounded-full border ${run.delta.includes('+') ? 'bg-red-500/10 text-red-500 border-red-500/20' : 'bg-green-500/10 text-green-500 border-green-500/20'
                                                        }`}>
                                                        {run.delta}
                                                    </span>
                                                )}
                                            </div>
                                        </td>
                                        <td className="px-6 py-5">
                                            {run.status === 'complete' ? (
                                                <div className="flex items-center gap-2 text-green-500 text-xs font-bold uppercase tracking-widest">
                                                    <CheckCircle2 size={14} /> Complete
                                                </div>
                                            ) : (
                                                <div className="flex items-center gap-2 text-red-500 text-xs font-bold uppercase tracking-widest">
                                                    <AlertCircle size={14} /> Error
                                                </div>
                                            )}
                                        </td>
                                        <td className="px-6 py-5 text-right">
                                            <button className="text-slate-500 hover:text-white transition-colors">
                                                <Play size={18} />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </section>

                {/* System Activity */}
                <section className="bg-blue-600/5 border border-blue-600/10 p-8 rounded-3xl flex flex-col md:flex-row items-center justify-between gap-8">
                    <div className="flex items-center gap-6">
                        <div className="w-16 h-16 rounded-2xl bg-blue-600/20 flex items-center justify-center text-blue-500 shadow-xl shadow-blue-500/10">
                            <Activity size={32} className="animate-pulse" />
                        </div>
                        <div>
                            <h3 className="text-xl font-bold text-white mb-1">Global Scan Throughput</h3>
                            <p className="text-slate-400 text-sm max-w-sm">
                                Total scans executed today across all verified domains. System load is within normal parameters.
                            </p>
                        </div>
                    </div>
                    <div className="flex items-baseline gap-2">
                        <span className="text-5xl font-extrabold text-white">458</span>
                        <span className="text-sm font-bold text-blue-500 uppercase tracking-widest">Scans / 24H</span>
                    </div>
                </section>
            </div>
        </DashboardLayout>
    );
}

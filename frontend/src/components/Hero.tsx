import { Shield, Lock, Search, Zap } from 'lucide-react';
import Link from 'next/link';

export default function Hero() {
    return (
        <div className="relative pt-32 pb-20 overflow-hidden">
            {/* Background Glow */}
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[600px] bg-blue-500/10 blur-[120px] rounded-full -z-10"></div>

            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
                <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-xs font-bold uppercase tracking-wider mb-8 animate-fade-in">
                    <Zap size={14} className="fill-current" />
                    Container-Aware Security Monitoring
                </div>

                <h1 className="text-5xl md:text-7xl font-extrabold text-white tracking-tight mb-6 max-w-4xl mx-auto leading-[1.1]">
                    Visualize Your <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-blue-600">Attack Surface</span> In Real-Time
                </h1>

                <p className="text-xl text-slate-400 max-w-2xl mx-auto mb-10 leading-relaxed">
                    The continuous threat exposure management platform for modern cloud-native teams. Discover, verify, and monitor containers, Kubernetes, and exposed APIs with deterministic intelligence.
                </p>

                <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
                    <Link href="/register" className="w-full sm:w-auto bg-accent hover:bg-blue-600 text-white px-8 py-4 rounded-xl font-bold text-lg transition-all shadow-xl shadow-blue-500/20 active:scale-95">
                        Start Monitoring Free
                    </Link>
                    <Link href="/dashboard" className="w-full sm:w-auto bg-slate-800 hover:bg-slate-700 text-slate-200 px-8 py-4 rounded-xl font-bold text-lg transition-all border border-slate-700 active:scale-95">
                        View Live Demo
                    </Link>
                </div>

                {/* Mock Dashboard Preview */}
                <div className="relative max-w-5xl mx-auto rounded-2xl border border-slate-700 bg-slate-900/50 backdrop-blur-xl shadow-2xl overflow-hidden group">
                    <div className="flex items-center gap-2 px-4 py-3 border-b border-slate-800 bg-slate-900/80">
                        <div className="flex gap-1.5">
                            <div className="w-3 h-3 rounded-full bg-red-500/50"></div>
                            <div className="w-3 h-3 rounded-full bg-yellow-500/50"></div>
                            <div className="w-3 h-3 rounded-full bg-green-500/50"></div>
                        </div>
                        <div className="ml-4 bg-slate-800 px-3 py-1 rounded text-[10px] text-slate-500 font-mono">
                            https://cortex.security/dashboard
                        </div>
                    </div>
                    <div className="p-4 md:p-8 bg-slate-900/40">
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-left">
                            {[
                                { label: 'Total Assets', value: '1,248', color: 'text-blue-400', icon: Search },
                                { label: 'Exposed Services', value: '42', color: 'text-orange-400', icon: Lock },
                                { label: 'Risk Score', value: '84/100', color: 'text-green-400', icon: Shield },
                            ].map((stat, i) => (
                                <div key={i} className="bg-slate-800/50 p-6 rounded-xl border border-slate-700/50">
                                    <div className="flex items-center justify-between mb-2">
                                        <span className="text-xs font-bold text-slate-500 uppercase">{stat.label}</span>
                                        <stat.icon size={16} className="text-slate-500" />
                                    </div>
                                    <div className={`text-3xl font-bold ${stat.color}`}>{stat.value}</div>
                                </div>
                            ))}
                        </div>
                        <div className="mt-8 bg-slate-800/30 border border-slate-700/50 rounded-xl p-6 text-left">
                            <div className="flex items-center justify-between mb-4">
                                <span className="text-sm font-bold text-slate-300">Active Vulnerability Scan: prod-cluster-01</span>
                                <span className="text-[10px] bg-blue-500/20 text-blue-400 px-2 py-0.5 rounded-full font-bold">RUNNING</span>
                            </div>
                            <div className="w-full bg-slate-700 h-1.5 rounded-full overflow-hidden">
                                <div className="bg-blue-500 h-full w-[65%] animate-pulse"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

import { Shield, Zap, Target, BarChart3, Globe, Lock } from 'lucide-react';
import Link from 'next/link';

export default function Hero() {
    return (
        <section className="relative pt-40 pb-32 overflow-hidden bg-slate-950">
            {/* Ambient Background Elements */}
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[120%] h-[1000px] bg-[radial-gradient(circle_at_center,rgba(59,130,246,0.08)_0%,transparent_70%)] -z-10"></div>
            <div className="absolute top-[-10%] right-[-10%] w-[40%] h-[40%] bg-blue-600/5 blur-[120px] rounded-full -z-10 animate-pulse"></div>

            <div className="max-w-7xl mx-auto px-6 text-center">
                {/* Badge */}
                <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/5 border border-blue-500/10 text-blue-400 text-[10px] font-black uppercase tracking-[0.2em] mb-12 shadow-inner">
                    <span className="relative flex h-2 w-2">
                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                        <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
                    </span>
                    Next-Gen Attack Surface Intelligence
                </div>

                {/* Main Headline */}
                <h1 className="text-6xl md:text-8xl font-black text-white tracking-tighter mb-8 max-w-5xl mx-auto leading-[0.95] uppercase italic">
                    Real-Time <br />
                    <span className="text-transparent bg-clip-text bg-gradient-to-b from-white to-slate-500">Threat Exposure</span>
                </h1>

                {/* Subtext */}
                <p className="text-lg md:text-xl text-slate-500 max-w-2xl mx-auto mb-14 leading-relaxed font-medium">
                    Continuous container-aware monitoring for the modern cloud stack.
                    Discover orphaned APIs, exposed Kubelets, and configuration drift before attackers do.
                </p>

                {/* CTAs */}
                <div className="flex flex-col sm:flex-row items-center justify-center gap-6 mb-24">
                    <Link href="/register" className="group relative px-10 py-5 bg-blue-600 hover:bg-blue-500 text-white rounded-2xl font-black text-sm uppercase tracking-widest transition-all shadow-2xl shadow-blue-500/20 active:scale-95 overflow-hidden">
                        <span className="relative z-10">Start Monitoring Free</span>
                        <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-1000"></div>
                    </Link>
                    <Link href="/demo" className="px-10 py-5 bg-slate-900 hover:bg-slate-800 text-slate-300 rounded-2xl font-black text-sm uppercase tracking-widest transition-all border border-slate-800 active:scale-95 shadow-xl">
                        Watch Demo
                    </Link>
                </div>

                {/* Dashboard Mockup Redesign */}
                <div className="relative max-w-6xl mx-auto rounded-[2.5rem] border border-white/5 bg-slate-900/20 backdrop-blur-3xl p-4 shadow-[0_0_100px_rgba(0,0,0,0.5)] group">
                    <div className="absolute inset-0 bg-blue-500/5 rounded-[2.5rem] -z-10 group-hover:bg-blue-500/10 transition-colors"></div>

                    <div className="bg-slate-950 rounded-[2rem] border border-white/10 overflow-hidden shadow-2xl">
                        {/* Browser Top Bar */}
                        <div className="flex items-center justify-between px-6 py-4 bg-slate-900/50 border-b border-white/5">
                            <div className="flex gap-2">
                                <div className="w-3 h-3 rounded-full bg-slate-800"></div>
                                <div className="w-3 h-3 rounded-full bg-slate-800"></div>
                                <div className="w-3 h-3 rounded-full bg-slate-800"></div>
                            </div>
                            <div className="flex items-center gap-3 bg-slate-950/80 px-4 py-1.5 rounded-full border border-white/5 text-[10px] text-slate-500 font-mono tracking-wider">
                                <Lock size={10} />
                                console.cortex.security
                            </div>
                            <div className="w-10"></div>
                        </div>

                        {/* Content Area */}
                        <div className="p-8 md:p-12">
                            <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
                                {[
                                    { label: 'Discovery Coverage', value: '99.4%', icon: Globe, color: 'text-blue-400' },
                                    { label: 'Active Targets', value: '14,208', icon: Target, color: 'text-white' },
                                    { label: 'Critical Risks', value: '03', icon: Shield, color: 'text-red-500' },
                                    { label: 'System Health', value: 'Nominal', icon: Zap, color: 'text-green-500' },
                                ].map((item, i) => (
                                    <div key={i} className="text-left">
                                        <div className="flex items-center gap-2 mb-2">
                                            <item.icon size={14} className="text-slate-600" />
                                            <span className="text-[10px] font-black text-slate-600 uppercase tracking-widest">{item.label}</span>
                                        </div>
                                        <div className={`text-3xl font-black ${item.color}`}>{item.value}</div>
                                    </div>
                                ))}
                            </div>

                            {/* Visual Chart Placeholder */}
                            <div className="h-64 w-full bg-slate-900/50 rounded-3xl border border-white/5 flex items-end justify-between p-8 gap-4 overflow-hidden relative">
                                {[40, 70, 45, 90, 65, 80, 55, 100, 75, 85, 60, 95].map((h, i) => (
                                    <div
                                        key={i}
                                        className="w-full bg-gradient-to-t from-blue-600/20 to-blue-500/40 rounded-t-lg transition-all duration-1000 group-hover:from-blue-600 group-hover:to-blue-400"
                                        style={{ height: `${h}%` }}
                                    ></div>
                                ))}
                                <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,transparent_0%,rgba(15,23,42,0.8)_100%)]"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
}

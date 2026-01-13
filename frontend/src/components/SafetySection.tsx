import { Gavel, Search, ShieldCheck } from 'lucide-react';

export default function SafetySection() {
    return (
        <section className="py-32 bg-slate-900/20 border-y border-white/5 relative">
            <div className="max-w-7xl mx-auto px-6">
                <div className="flex flex-col lg:flex-row items-center justify-between gap-16">
                    <div className="max-w-2xl text-left">
                        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/5 border border-emerald-500/10 text-emerald-500 text-[10px] font-black uppercase tracking-[0.2em] mb-8">
                            Non-Destructive Scanning
                        </div>
                        <h2 className="text-4xl md:text-6xl font-black text-white tracking-tighter uppercase italic mb-8">
                            Legally Safe & <br />
                            <span className="text-slate-500">Security-First</span>
                        </h2>
                        <p className="text-lg text-slate-500 font-medium leading-relaxed mb-10">
                            Cortex is designed for professional security teams who value target stability and ethical boundaries.
                            Our engine strictly follows a "no-exploit" policy, focusing purely on exposure intelligence
                            without risking infrastructure uptime.
                        </p>

                        <div className="space-y-6">
                            {[
                                { title: "Strict Scope Enforcement", icon: Gavel, desc: "Scanning is only permitted for verified domains with valid ownership records." },
                                { title: "Passive-First Methodology", icon: Search, desc: "We prioritize non-invasive techniques that never trigger stability warnings." },
                                { title: "SOC2 & GDPR Aligned", icon: ShieldCheck, desc: "Built with enterprise data privacy and data handling standards at the core." }
                            ].map((item, i) => (
                                <div key={i} className="flex gap-4">
                                    <div className="mt-1">
                                        <item.icon className="text-blue-500 w-5 h-5" />
                                    </div>
                                    <div>
                                        <h4 className="text-sm font-black text-white uppercase tracking-widest mb-1">{item.title}</h4>
                                        <p className="text-xs text-slate-500 font-medium">{item.desc}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="relative w-full lg:w-96 aspect-square bg-slate-900 border border-white/5 rounded-[3rem] p-12 flex flex-col justify-center shadow-2xl overflow-hidden group">
                        <div className="absolute top-0 right-0 w-32 h-32 bg-blue-500/10 blur-3xl -z-10 group-hover:bg-blue-500/20 transition-colors"></div>
                        <div className="text-center">
                            <div className="w-20 h-20 bg-blue-600 rounded-3xl mx-auto mb-8 flex items-center justify-center shadow-lg shadow-blue-500/20">
                                <ShieldCheck className="text-white w-10 h-10" />
                            </div>
                            <div className="text-2xl font-black text-white uppercase italic tracking-tighter mb-2">Verified Scoped</div>
                            <div className="text-[10px] font-black text-blue-500 uppercase tracking-widest">Safe to Execute</div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
}

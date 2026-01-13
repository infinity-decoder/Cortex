import { Boxes, Cpu, Unlink, Eye, AlertTriangle, ShieldCheck } from 'lucide-react';

const pillars = [
    {
        title: "Container Probing",
        desc: "Deep inspection of Kubelets and Docker APIs to identify unauthenticated control plane exposure.",
        icon: Boxes,
        color: "bg-blue-500",
        tag: "Infrastucture"
    },
    {
        title: "API Surface Mapping",
        desc: "Automatic discovery of orphaned endpoints and deprecated API versions leaking sensitive data.",
        icon: Cpu,
        color: "bg-purple-500",
        tag: "Application"
    },
    {
        title: "Exposure Chaining",
        desc: "Visualize how minor misconfigurations lead to full cluster compromise through attack path analysis.",
        icon: Unlink,
        color: "bg-orange-500",
        tag: "Intelligence"
    },
    {
        title: "Continuous Observability",
        desc: "24/7 monitoring of DNS, TLS, and Service headers with real-time delta alerting.",
        icon: Eye,
        color: "bg-emerald-500",
        tag: "Monitoring"
    },
    {
        title: "Risk Quantification",
        desc: "Deterministic risk scoring based on actual exploitability, not just CVSS noise.",
        icon: AlertTriangle,
        color: "bg-red-500",
        tag: "Analytics"
    },
    {
        title: "Verified Compliance",
        desc: "Ensuring all discovered assets are within your legal scope through multi-step verification.",
        icon: ShieldCheck,
        color: "bg-indigo-500",
        tag: "Safety"
    }
];

export default function Features() {
    return (
        <section id="features" className="py-32 bg-slate-950 relative overflow-hidden">
            <div className="max-w-7xl mx-auto px-6">
                <div className="text-left mb-24 max-w-3xl">
                    <h2 className="text-4xl md:text-6xl font-black text-white tracking-tighter uppercase italic mb-8">
                        The Core Pillars of <br />
                        <span className="text-slate-500">Continuous Security</span>
                    </h2>
                    <p className="text-lg text-slate-500 font-medium leading-relaxed">
                        Modern security teams need more than just a port scanner. Cortex provides the context
                        required to understand your true exposure in a containerized world.
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    {pillars.map((pillar, i) => (
                        <div key={i} className="group relative bg-slate-900/30 border border-white/5 rounded-[2.5rem] p-10 hover:bg-slate-900/50 transition-all duration-500 cursor-default">
                            <div className={`w-14 h-14 ${pillar.color}/10 rounded-2xl flex items-center justify-center mb-8 group-hover:scale-110 transition-transform`}>
                                <pillar.icon className={`${pillar.color.replace('bg-', 'text-')} w-8 h-8`} />
                            </div>

                            <div className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-600 mb-2">{pillar.tag}</div>
                            <h3 className="text-2xl font-black text-white mb-4 tracking-tight">{pillar.title}</h3>
                            <p className="text-slate-500 text-sm leading-relaxed font-medium">
                                {pillar.desc}
                            </p>

                            <div className="absolute top-6 right-8 opacity-0 group-hover:opacity-100 transition-opacity">
                                <div className="w-2 h-2 rounded-full bg-blue-500 animate-pulse"></div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Background Decorative Blur */}
            <div className="absolute bottom-0 right-0 w-[500px] h-[500px] bg-blue-600/5 blur-[120px] rounded-full -z-10"></div>
        </section>
    );
}

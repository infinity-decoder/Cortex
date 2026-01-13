import { UserPlus, Search, ShieldCheck, Bell } from 'lucide-react';

const steps = [
    {
        title: "Identify & Verify",
        desc: "Connect your root domains. We perform a non-invasive ownership check to ensure you only scan what you own.",
        icon: UserPlus,
    },
    {
        title: "Deep Discovery",
        desc: "Cortex maps your entire surface, finding orphaned subdomains and hidden containerized endpoints.",
        icon: Search,
    },
    {
        title: "Risk Analysis",
        desc: "Our engine finger-prints every service and probes for unauthenticated APIs and misconfigurations.",
        icon: ShieldCheck,
    },
    {
        title: "Delta Alerting",
        desc: "Stay notified of new exposures. Get actionable remediation steps the moment a drift is detected.",
        icon: Bell,
    }
];

export default function HowItWorks() {
    return (
        <section className="py-32 bg-slate-950 relative">
            <div className="max-w-7xl mx-auto px-6">
                <div className="text-center mb-24">
                    <h2 className="text-4xl md:text-6xl font-black text-white tracking-tighter uppercase italic mb-8">
                        Deterministic <span className="text-slate-500">Execution</span>
                    </h2>
                    <p className="text-lg text-slate-500 font-medium max-w-2xl mx-auto">
                        From initial scope definition to real-time risk remediation in four clean steps.
                    </p>
                </div>

                <div className="relative">
                    {/* Connecting Line (Desktop) */}
                    <div className="hidden lg:block absolute top-1/2 left-0 w-full h-px bg-white/5 -translate-y-1/2 -z-10"></div>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 text-left">
                        {steps.map((step, i) => (
                            <div key={i} className="flex flex-col">
                                <div className="w-16 h-16 bg-slate-900 border border-white/5 rounded-2xl flex items-center justify-center text-blue-500 mb-8 shadow-xl">
                                    <step.icon size={28} />
                                </div>
                                <h3 className="text-xl font-black text-white mb-4 tracking-tight">0{i + 1}. {step.title}</h3>
                                <p className="text-slate-500 text-sm leading-relaxed font-medium">
                                    {step.desc}
                                </p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </section>
    );
}

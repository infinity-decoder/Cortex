import { Check } from 'lucide-react';
import Link from 'next/link';

const tiers = [
    {
        name: "Discovery",
        price: "Free",
        desc: "Essential reconnaissance for individuals and open-source projects.",
        features: [
            "Recursive DNS Enumeration",
            "3 Verified Domains",
            "Port & Service Discovery",
            "Weekly Scan Frequency",
            "Community Support"
        ],
        cta: "Start Free",
        highlight: false
    },
    {
        name: "Professional",
        price: "$490",
        period: "/mo",
        desc: "Full attack surface management for growing product teams.",
        features: [
            "Unlimited Verified Domains",
            "Continuous Delta Monitoring",
            "Kubernetes & Docker Probing",
            "API Leak Detection",
            "Priority Slack Support"
        ],
        cta: "Start 14-Day Trial",
        highlight: true
    },
    {
        name: "Enterprise",
        price: "Custom",
        desc: "Global-scale vulnerability intelligence for security operations.",
        features: [
            "Multi-Org Resource Isolation",
            "SAML / SSO Integration",
            "Attack Path Visualization",
            "Compliance Audit Reports",
            "Designated Security Lead"
        ],
        cta: "Contact Sales",
        highlight: false
    }
];

export default function Pricing() {
    return (
        <section id="pricing" className="py-32 bg-slate-950">
            <div className="max-w-7xl mx-auto px-6">
                <div className="text-center mb-24">
                    <h2 className="text-4xl md:text-6xl font-black text-white tracking-tighter uppercase italic mb-8">
                        Scale Your <span className="text-slate-500">Security Coverage</span>
                    </h2>
                    <p className="text-lg text-slate-500 font-medium max-w-2xl mx-auto leading-relaxed">
                        Deterministic pricing with no hidden seat costs. Choose the tier that
                        matches your infrastructure's complexity.
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8 items-stretch">
                    {tiers.map((tier, i) => (
                        <div
                            key={i}
                            className={`relative rounded-[2.5rem] p-10 flex flex-col transition-all duration-500 ${tier.highlight
                                    ? 'bg-blue-600 shadow-[0_0_80px_rgba(37,99,235,0.15)] border-white/20'
                                    : 'bg-slate-900/30 border border-white/5 hover:bg-slate-900/50'
                                }`}
                        >
                            {tier.highlight && (
                                <div className="absolute -top-4 left-1/2 -translate-x-1/2 bg-white text-blue-600 px-4 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest shadow-xl">
                                    Most Popular
                                </div>
                            )}

                            <div className="mb-10">
                                <h3 className={`text-xl font-black uppercase tracking-widest mb-4 ${tier.highlight ? 'text-white' : 'text-slate-400'}`}>
                                    {tier.name}
                                </h3>
                                <div className="flex items-baseline gap-1 mb-4">
                                    <span className="text-5xl font-black text-white tracking-tighter">{tier.price}</span>
                                    {tier.period && <span className={`text-sm font-bold ${tier.highlight ? 'text-blue-100' : 'text-slate-600'}`}>{tier.period}</span>}
                                </div>
                                <p className={`text-sm font-medium leading-relaxed ${tier.highlight ? 'text-blue-100' : 'text-slate-500'}`}>
                                    {tier.desc}
                                </p>
                            </div>

                            <div className="space-y-4 mb-12 flex-grow">
                                {tier.features.map((feature, j) => (
                                    <div key={j} className="flex items-center gap-3">
                                        <div className={`w-5 h-5 rounded-full flex items-center justify-center ${tier.highlight ? 'bg-white/10' : 'bg-slate-800'}`}>
                                            <Check size={12} className={tier.highlight ? 'text-white' : 'text-blue-500'} />
                                        </div>
                                        <span className={`text-xs font-bold ${tier.highlight ? 'text-white' : 'text-slate-400'}`}>{feature}</span>
                                    </div>
                                ))}
                            </div>

                            <Link
                                href="/register"
                                className={`block w-full py-4 rounded-2xl text-center text-xs font-black uppercase tracking-widest transition-all active:scale-95 ${tier.highlight
                                        ? 'bg-white text-blue-600 hover:bg-slate-100 shadow-xl shadow-blue-900/20'
                                        : 'bg-slate-800 text-white hover:bg-slate-700'
                                    }`}
                            >
                                {tier.cta}
                            </Link>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}

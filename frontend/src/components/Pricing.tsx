import { Check } from 'lucide-react';
import Link from 'next/link';

const plans = [
    {
        name: 'Free',
        price: '$0',
        description: 'Perfect for individual developers and small side projects.',
        features: ['1 Root Domain', 'Unlimited Subdomains', 'Daily Scans', 'Basic Findings Detail', 'Community Support'],
        cta: 'Get Started',
        premium: false,
    },
    {
        name: 'Pro',
        price: '$49',
        description: 'Ideal for growing SaaS startups and DevOps teams.',
        features: ['5 Root Domains', 'Hourly Scans', 'Priority Email Alerts', 'Remediation Guidance', 'Team RBAC (3 Users)'],
        cta: 'Start Pro Trial',
        premium: true,
    },
    {
        name: 'Enterprise',
        price: 'Custom',
        description: 'For organizations requires global attack surface scale.',
        features: ['Unlimited Domains', 'Real-time Monitoring', 'Webhook Integrations', 'Dedicated Account Manager', 'Custom Data Retention'],
        cta: 'Contact Sales',
        premium: false,
    },
];

export default function Pricing() {
    return (
        <section id="pricing" className="py-24 relative overflow-hidden">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="text-center mb-16">
                    <h2 className="text-3xl md:text-5xl font-bold text-white mb-4">Simple, Predictable Pricing</h2>
                    <p className="text-slate-400 max-w-2xl mx-auto">
                        Scale your monitoring as your infrastructure grows. Start free, no credit card required.
                    </p>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {plans.map((plan, i) => (
                        <div key={i} className={`p-8 rounded-3xl border ${plan.premium ? 'bg-slate-800 border-blue-500 shadow-2xl shadow-blue-500/10' : 'bg-slate-800/40 border-slate-700/50'} flex flex-col relative`}>
                            {plan.premium && (
                                <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-blue-500 text-[10px] font-extrabold px-3 py-1 rounded-full uppercase tracking-widest">
                                    Most Popular
                                </div>
                            )}
                            <h3 className="text-xl font-bold text-white mb-2">{plan.name}</h3>
                            <div className="flex items-baseline gap-1 mb-4">
                                <span className="text-4xl font-extrabold text-white">{plan.price}</span>
                                {plan.price !== 'Custom' && <span className="text-slate-500 font-medium">/mo</span>}
                            </div>
                            <p className="text-sm text-slate-400 mb-8 leading-relaxed">
                                {plan.description}
                            </p>

                            <div className="space-y-4 mb-10 flex-1">
                                {plan.features.map((feature, j) => (
                                    <div key={j} className="flex items-center gap-3 text-sm text-slate-300">
                                        <Check size={16} className="text-blue-500 flex-shrink-0" />
                                        {feature}
                                    </div>
                                ))}
                            </div>

                            <Link
                                href={plan.name === 'Enterprise' ? '/contact' : '/register'}
                                className={`w-full py-3 rounded-xl font-bold text-center transition-all active:scale-95 ${plan.premium ? 'bg-blue-600 hover:bg-blue-500 text-white shadow-lg' : 'bg-slate-700 hover:bg-slate-600 text-slate-200'}`}
                            >
                                {plan.cta}
                            </Link>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}

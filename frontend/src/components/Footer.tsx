import Link from 'next/link';
import { Shield, Github, Twitter, Linkedin, Mail } from 'lucide-react';

const links = {
    Product: [
        { name: 'Features', href: '#features' },
        { name: 'Intelligence', href: '#intelligence' },
        { name: 'Scanner', href: '/dashboard' },
        { name: 'API Docs', href: '/docs' },
    ],
    Legal: [
        { name: 'Privacy Policy', href: '/privacy' },
        { name: 'Terms of Service', href: '/terms' },
        { name: 'Scan Policy', href: '/scan-policy' },
        { name: 'Scope Agreement', href: '/scope' },
    ],
    Resources: [
        { name: 'Security Blog', href: '/blog' },
        { name: 'CVE Library', href: '/cve' },
        { name: 'Status', href: '/status' },
        { name: 'Newsletter', href: '/newsletter' },
    ]
};

export default function Footer() {
    return (
        <footer className="bg-slate-950 border-t border-white/5 pt-24 pb-12">
            <div className="max-w-7xl mx-auto px-6">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-12 mb-20">
                    {/* Brand Column */}
                    <div className="lg:col-span-2 space-y-8 text-left">
                        <Link href="/" className="flex items-center gap-3">
                            <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/20">
                                <Shield className="text-white w-6 h-6" />
                            </div>
                            <span className="text-2xl font-black text-white tracking-tighter uppercase italic">Cortex</span>
                        </Link>
                        <p className="text-slate-500 text-sm font-medium leading-relaxed max-w-sm">
                            Deterministic threat exposure management for the cloud-native era.
                            Stay secure with legally-safe, continuous intelligence.
                        </p>
                        <div className="flex items-center gap-5">
                            {[Github, Twitter, Linkedin, Mail].map((Icon, i) => (
                                <a key={i} href="#" className="text-slate-600 hover:text-white transition-colors">
                                    <Icon size={20} />
                                </a>
                            ))}
                        </div>
                    </div>

                    {/* Link Columns */}
                    {Object.entries(links).map(([title, items], i) => (
                        <div key={i} className="text-left">
                            <h4 className="text-[10px] font-black uppercase tracking-[0.2em] text-white mb-6">
                                {title}
                            </h4>
                            <ul className="space-y-4">
                                {items.map((link, j) => (
                                    <li key={j}>
                                        <Link
                                            href={link.href}
                                            className="text-sm font-bold text-slate-500 hover:text-blue-400 transition-colors"
                                        >
                                            {link.name}
                                        </Link>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>

                <div className="pt-12 border-t border-white/5 flex flex-col md:flex-row items-center justify-between gap-6">
                    <p className="text-[10px] font-black text-slate-600 uppercase tracking-widest">
                        © 2026 Cortex Security. All rights reserved.
                    </p>
                    <div className="flex items-center gap-8">
                        <div className="flex items-center gap-2">
                            <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
                            <span className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Global Scan Status: Optimal</span>
                        </div>
                        <span className="text-[10px] font-black text-slate-700 uppercase tracking-widest">v0.4.2-PROD</span>
                    </div>
                </div>
            </div>
        </footer>
    );
}

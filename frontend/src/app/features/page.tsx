import Navbar from '@/components/Navbar';
import Footer from '@/components/Footer';
import { Boxes, Cpu, Unlink, Eye, AlertTriangle, ShieldCheck, ChevronRight } from 'lucide-react';
import Link from 'next/link';

const features = [
    {
        title: "Recursive DNS & Subdomain Enumeration",
        desc: "Our engine uses a multi-source approach, combining passive DNS records with brute-force discovery using custom wordlists tailored for cloud-native infrastructure.",
        icon: Boxes,
        details: ["Wildcard detection", "CNAME flattened mapping", "Historical record analysis"]
    },
    {
        title: "K8s & Docker Control Plane Probing",
        desc: "Automatically identify unauthenticated Kubelets, Docker Remote APIs, and Etcd endpoints that could lead to full cluster takeovers.",
        icon: Cpu,
        details: ["Version fingerprinting", "RBAC policy leakage", "Container registry analysis"]
    },
    {
        title: "Dynamic Attack Surface Mapping",
        desc: "Stay ahead of shadow IT. Cortex continuously scans your verified IP ranges and domains to detect new assets as they are provisioned.",
        icon: Unlink,
        details: ["Cloud provider integration", "Delta reporting", "Orphaned subdomain tracking"]
    }
];

export default function FeaturesPage() {
    return (
        <div className="min-h-screen bg-slate-950 text-slate-200">
            <Navbar />

            <main className="pt-40 pb-32">
                <div className="max-w-7xl mx-auto px-6">
                    <div className="text-center mb-24">
                        <h1 className="text-5xl md:text-7xl font-black text-white tracking-tighter uppercase italic mb-8">
                            Technical <span className="text-slate-500">Capability</span>
                        </h1>
                        <p className="text-lg text-slate-500 font-medium max-w-2xl mx-auto leading-relaxed">
                            A deep dive into the methodology and engineering behind the Cortex exposure engine.
                        </p>
                    </div>

                    <div className="space-y-12">
                        {features.map((f, i) => (
                            <div key={i} className="bg-slate-900/40 border border-white/5 rounded-[3rem] p-12 lg:p-20 flex flex-col lg:flex-row gap-16 items-center">
                                <div className="flex-1 text-left">
                                    <div className="w-16 h-16 bg-blue-500/10 rounded-2xl flex items-center justify-center text-blue-500 mb-10">
                                        <f.icon size={32} />
                                    </div>
                                    <h2 className="text-3xl md:text-5xl font-black text-white tracking-tighter uppercase italic mb-8 leading-[0.9]">{f.title}</h2>
                                    <p className="text-lg text-slate-500 font-medium leading-relaxed mb-10">{f.desc}</p>

                                    <ul className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        {f.details.map((detail, j) => (
                                            <li key={j} className="flex items-center gap-3 text-xs font-black uppercase tracking-widest text-slate-400">
                                                <div className="w-1.5 h-1.5 rounded-full bg-blue-500"></div>
                                                {detail}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                                <div className="flex-1 w-full aspect-video bg-slate-950 rounded-[2rem] border border-white/10 shadow-2xl relative overflow-hidden flex items-center justify-center group">
                                    <div className="absolute inset-0 bg-blue-500/5 -z-10 group-hover:bg-blue-500/10 transition-colors"></div>
                                    <div className="text-[10px] font-mono font-bold text-slate-700 tracking-widest uppercase">Visual Evidence Placeholder</div>
                                </div>
                            </div>
                        ))}
                    </div>

                    {/* CTA */}
                    <div className="mt-32 text-center">
                        <Link href="/register" className="inline-flex items-center gap-4 bg-blue-600 hover:bg-blue-500 text-white px-12 py-6 rounded-3xl font-black text-sm uppercase tracking-widest transition-all shadow-2xl shadow-blue-500/20 active:scale-95">
                            Start Free Discovery <ChevronRight size={20} />
                        </Link>
                    </div>
                </div>
            </main>

            <Footer />
        </div>
    );
}

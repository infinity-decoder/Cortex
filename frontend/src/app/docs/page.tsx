import Navbar from '@/components/Navbar';
import Footer from '@/components/Footer';
import { Book, Code, Terminal, Shield, Search, Zap, ChevronRight } from 'lucide-react';
import Link from 'next/link';

const categories = [
    { title: "Getting Started", icon: Book, desc: "Quick start guide to setting up your first verified domain.", slug: "welcome" },
    { title: "Discovery Engine", icon: Search, desc: "How our DNS and passive discovery mechanisms work.", slug: "discovery" },
    { title: "Scan Methodologies", icon: Zap, desc: "Details on non-invasive vs deep container probing.", slug: "scans" },
    { title: "API Documentation", icon: Code, desc: "Reference for our REST API and CLI integrations.", slug: "api" },
    { title: "Compliance & Safety", icon: Shield, desc: "Ownership verification and our legally-safe scanning policy.", slug: "safety" },
    { title: "CLI Reference", icon: Terminal, desc: "Command line tools for automated CI/CD security workflows.", slug: "cli" },
];

export default function Docs() {
    return (
        <div className="min-h-screen bg-slate-950 text-slate-200">
            <Navbar />

            <main className="pt-40 pb-32">
                <div className="max-w-7xl mx-auto px-6">
                    <div className="text-center mb-24">
                        <h1 className="text-5xl md:text-7xl font-black text-white tracking-tighter uppercase italic mb-8">
                            Knowledge <span className="text-slate-500">Center</span>
                        </h1>
                        <p className="text-lg text-slate-500 font-medium max-w-2xl mx-auto leading-relaxed">
                            Everything you need to master continuous threat exposure management with Cortex.
                        </p>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mb-24">
                        {categories.map((cat) => (
                            <Link
                                href={`/docs/${cat.slug}`}
                                key={cat.slug}
                                className="group bg-slate-900/40 border border-white/5 p-10 rounded-[2.5rem] hover:bg-slate-900/60 transition-all duration-300"
                            >
                                <div className="w-12 h-12 bg-blue-500/10 rounded-2xl flex items-center justify-center text-blue-500 mb-8 font-black group-hover:scale-110 transition-transform">
                                    <cat.icon size={24} />
                                </div>
                                <h3 className="text-xl font-black text-white mb-4 uppercase tracking-tight">{cat.title}</h3>
                                <p className="text-sm text-slate-500 font-medium leading-relaxed mb-6">
                                    {cat.desc}
                                </p>
                                <div className="flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-blue-500 group-hover:gap-3 transition-all">
                                    Explore Topic <ChevronRight size={14} />
                                </div>
                            </Link>
                        ))}
                    </div>

                    {/* Quick Search / Help Bar */}
                    <div className="bg-blue-600/5 border border-blue-600/10 p-12 rounded-[3rem] text-center max-w-4xl mx-auto">
                        <h4 className="text-2xl font-black text-white uppercase italic tracking-tighter mb-4">Can't find what you're looking for?</h4>
                        <p className="text-slate-500 text-sm font-medium mb-8"> Our security engineers are standing by to help with your implementation.</p>
                        <div className="flex justify-center gap-6">
                            <Link href="/contact" className="bg-blue-600 hover:bg-blue-500 text-white px-8 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all">
                                Contact Support
                            </Link>
                            <Link href="https://github.com/cortex-security/cortex" className="bg-slate-900 hover:bg-slate-800 text-white px-8 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest border border-slate-800 transition-all">
                                Github Discussions
                            </Link>
                        </div>
                    </div>
                </div>
            </main>

            <Footer />
        </div>
    );
}

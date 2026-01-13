import Link from 'next/link';
import { Shield } from 'lucide-react';

export default function Navbar() {
    return (
        <nav className="fixed top-0 w-full z-50 transition-all duration-300 border-b border-white/5 bg-slate-950/80 backdrop-blur-xl">
            <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
                <Link href="/" className="flex items-center gap-3 group">
                    <div className="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/20 group-hover:scale-110 transition-transform">
                        <Shield className="text-white w-6 h-6" />
                    </div>
                    <span className="text-2xl font-black text-white tracking-tighter uppercase italic">Cortex</span>
                </Link>

                <div className="hidden lg:flex items-center gap-10">
                    {['Features', 'Intelligence', 'Pricing', 'Security'].map((item) => (
                        <Link
                            key={item}
                            href={`#${item.toLowerCase()}`}
                            className="text-sm font-bold text-slate-400 hover:text-white transition-colors uppercase tracking-widest"
                        >
                            {item}
                        </Link>
                    ))}
                </div>

                <div className="flex items-center gap-6">
                    <Link href="/login" className="text-sm font-bold text-slate-300 hover:text-white transition-colors uppercase tracking-widest">
                        Login
                    </Link>
                    <Link href="/register" className="bg-white text-slate-950 hover:bg-slate-200 px-6 py-2.5 rounded-full text-xs font-black uppercase tracking-widest transition-all shadow-xl active:scale-95">
                        Get Started
                    </Link>
                </div>
            </div>
        </nav>
    );
}

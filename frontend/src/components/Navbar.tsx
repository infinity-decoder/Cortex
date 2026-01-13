import Link from 'next/link';
import Image from 'next/image';

export default function Navbar() {
    return (
        <nav className="fixed top-0 w-full z-50 bg-background/80 backdrop-blur-md border-b border-slate-800">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                    <div className="flex items-center gap-2">
                        <Link href="/" className="flex items-center gap-2">
                            <Image
                                src="/cortex_SaaS.png"
                                alt="Cortex Logo"
                                width={32}
                                height={32}
                                className="rounded"
                            />
                            <span className="text-xl font-bold tracking-tight text-white">Cortex</span>
                        </Link>
                    </div>

                    <div className="hidden md:flex items-center gap-8">
                        <Link href="#features" className="text-sm text-slate-400 hover:text-white transition-colors">Features</Link>
                        <Link href="#pricing" className="text-sm text-slate-400 hover:text-white transition-colors">Pricing</Link>
                        <Link href="/docs" className="text-sm text-slate-400 hover:text-white transition-colors">Documentation</Link>
                    </div>

                    <div className="flex items-center gap-4">
                        <Link href="/login" className="text-sm font-medium text-slate-300 hover:text-white transition-colors">
                            Login
                        </Link>
                        <Link href="/register" className="bg-accent hover:bg-blue-600 px-4 py-2 rounded-lg text-sm font-bold transition-all transform active:scale-95 shadow-lg shadow-blue-500/20">
                            Get Started
                        </Link>
                    </div>
                </div>
            </div>
        </nav>
    );
}

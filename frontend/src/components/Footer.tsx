import Image from 'next/image';
import Link from 'next/link';

export default function Footer() {
    return (
        <footer className="bg-slate-900 border-t border-slate-800 py-16">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-12 lg:gap-8 mb-12">
                    <div className="col-span-2">
                        <Link href="/" className="flex items-center gap-2 mb-6">
                            <Image
                                src="/cortex_SaaS.png"
                                alt="Cortex Logo"
                                width={32}
                                height={32}
                                className="rounded grayscale hover:grayscale-0 transition-all"
                            />
                            <span className="text-xl font-bold tracking-tight text-white">Cortex</span>
                        </Link>
                        <p className="text-slate-500 text-sm max-w-xs leading-relaxed">
                            Professional-grade attack surface monitoring for cloud-native infrastructure. deterministic, legally-safe, and secure.
                        </p>
                    </div>

                    <div>
                        <h4 className="text-white font-bold text-sm mb-6 uppercase tracking-wider">Product</h4>
                        <ul className="space-y-4 text-sm">
                            <li><Link href="#features" className="text-slate-500 hover:text-white transition-colors">Features</Link></li>
                            <li><Link href="/dashboard" className="text-slate-500 hover:text-white transition-colors">Scanner</Link></li>
                            <li><Link href="#pricing" className="text-slate-500 hover:text-white transition-colors">Pricing</Link></li>
                            <li><Link href="/docs" className="text-slate-500 hover:text-white transition-colors">API Docs</Link></li>
                        </ul>
                    </div>

                    <div>
                        <h4 className="text-white font-bold text-sm mb-6 uppercase tracking-wider">Legal</h4>
                        <ul className="space-y-4 text-sm">
                            <li><Link href="/privacy" className="text-slate-500 hover:text-white transition-colors">Privacy Policy</Link></li>
                            <li><Link href="/terms" className="text-slate-500 hover:text-white transition-colors">Terms of Service</Link></li>
                            <li><Link href="/scan-policy" className="text-slate-500 hover:text-white transition-colors">Scan Policy</Link></li>
                            <li><Link href="/compliance" className="text-slate-500 hover:text-white transition-colors">Compliance</Link></li>
                        </ul>
                    </div>

                    <div>
                        <h4 className="text-white font-bold text-sm mb-6 uppercase tracking-wider">Connect</h4>
                        <ul className="space-y-4 text-sm">
                            <li><Link href="https://github.com/infinity-decoder/Cortex" target="_blank" className="text-slate-500 hover:text-white transition-colors">GitHub</Link></li>
                            <li><Link href="/blog" className="text-slate-500 hover:text-white transition-colors">Security Blog</Link></li>
                            <li><Link href="/newsletter" className="text-slate-500 hover:text-white transition-colors">Newsletter</Link></li>
                            <li><Link href="#contact" className="text-slate-500 hover:text-white transition-colors">Contact</Link></li>
                        </ul>
                    </div>
                </div>

                <div className="pt-8 border-t border-slate-800 flex flex-col md:row items-center justify-between gap-4">
                    <p className="text-slate-600 text-xs text-center md:text-left">
                        © {new Date().getFullYear()} Cortex Security. All rights reserved.
                        Designed for the 2026-2030 CTEM Evolution.
                    </p>
                    <div className="flex gap-6 text-slate-600 text-xs">
                        <span>Status: All Systems Operational</span>
                    </div>
                </div>
            </div>
        </footer>
    );
}

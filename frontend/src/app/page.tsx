import Navbar from '@/components/Navbar';
import Hero from '@/components/Hero';
import Features from '@/components/Features';
import Pricing from '@/components/Pricing';
import Footer from '@/components/Footer';
import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <Navbar />
      <main>
        <Hero />
        <Features />

        {/* Trust Section */}
        <section className="py-20 border-y border-slate-800 bg-slate-800/20">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex flex-col md:flex-row items-center justify-between gap-12">
              <div className="max-w-xl">
                <h3 className="text-2xl font-bold text-white mb-4">Legally Safe & Secure Recon</h3>
                <p className="text-slate-400 leading-relaxed">
                  Cortex enforces strict ownership verification and read-only probing. We never exploit, never brute force, and never compromise target stability. Designed for security teams that value trust.
                </p>
              </div>
              <div className="flex flex-wrap gap-8 items-center grayscale opacity-50 justify-center">
                <div className="text-white font-mono font-bold text-xl uppercase tracking-tighter shadow-sm">SOC2 Ready</div>
                <div className="text-white font-mono font-bold text-xl uppercase tracking-tighter shadow-sm">GDPR Compliant</div>
                <div className="text-white font-mono font-bold text-xl uppercase tracking-tighter shadow-sm">CTEM Aligned</div>
              </div>
            </div>
          </div>
        </section>

        <Pricing />

        {/* Final CTA */}
        <section className="py-24 relative overflow-hidden text-center">
          <div className="absolute inset-0 bg-blue-600/5 -z-10"></div>
          <div className="max-w-4xl mx-auto px-4">
            <h2 className="text-3xl md:text-5xl font-bold text-white mb-6 tracking-tight">Ready to secure your cloud surface?</h2>
            <p className="text-slate-400 text-lg mb-10 max-w-2xl mx-auto leading-relaxed">Join DevOps teams using Cortex for deterministic threat exposure management.</p>
            <div className="flex flex-col sm:flex-row justify-center gap-4">
              <Link href="/register" className="bg-blue-600 hover:bg-blue-500 text-white px-8 py-4 rounded-xl font-bold transition-all shadow-xl shadow-blue-500/20 active:scale-95">
                Create Free Organization
              </Link>
              <Link href="/contact" className="bg-slate-800 hover:bg-slate-700 text-white px-8 py-4 rounded-xl font-bold transition-all border border-slate-700 active:scale-95">
                Talk to Sales
              </Link>
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}

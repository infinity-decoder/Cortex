import Navbar from '@/components/Navbar';
import Hero from '@/components/Hero';
import SocialProof from '@/components/SocialProof';
import Features from '@/components/Features';
import HowItWorks from '@/components/HowItWorks';
import SafetySection from '@/components/SafetySection';
import Pricing from '@/components/Pricing';
import Footer from '@/components/Footer';
import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-200 selection:bg-blue-500/30">
      <Navbar />

      <main>
        <Hero />
        <SocialProof />

        <div id="features">
          <Features />
        </div>

        <HowItWorks />

        <SafetySection />

        <div id="pricing">
          <Pricing />
        </div>

        {/* Final CTA Section */}
        <section className="py-24 relative overflow-hidden bg-slate-950 border-t border-white/5">
          <div className="absolute inset-0 bg-blue-600/5 -z-10 blur-3xl"></div>
          <div className="max-w-4xl mx-auto px-6 text-center">
            <h2 className="text-5xl md:text-7xl font-black text-white mb-8 tracking-tighter uppercase italic">
              Ready to Secure Your <br />
              <span className="text-slate-500">Cloud Surface?</span>
            </h2>
            <p className="text-slate-500 text-lg mb-12 max-w-2xl mx-auto font-medium leading-relaxed">
              Join high-velocity teams using Cortex for deterministic, ethically-safe threat exposure management.
            </p>
            <div className="flex flex-col sm:flex-row justify-center gap-6">
              <Link href="/register" className="bg-blue-600 hover:bg-blue-500 text-white px-10 py-5 rounded-2xl font-black text-xs uppercase tracking-widest transition-all shadow-xl shadow-blue-500/20 active:scale-95">
                Create Free Account
              </Link>
              <Link href="/contact" className="bg-slate-900 hover:bg-slate-800 text-white px-10 py-5 rounded-2xl font-black text-xs uppercase tracking-widest transition-all border border-slate-800 active:scale-95">
                Review Scan Policy
              </Link>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}

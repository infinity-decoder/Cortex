export default function SocialProof() {
    const brands = ["Cloudflare", "Datadog", "Crowdstrike", "Snyk", "Vercel", "HashiCorp"];

    return (
        <section className="py-20 bg-slate-950 border-y border-white/5 relative overflow-hidden">
            <div className="max-w-7xl mx-auto px-6">
                <div className="text-center mb-12">
                    <p className="text-[10px] font-black text-slate-600 uppercase tracking-[0.3em]">
                        Trusted by High-Velocity Engineering Teams
                    </p>
                </div>

                <div className="flex flex-wrap justify-center items-center gap-x-16 gap-y-10 grayscale opacity-40 hover:opacity-100 transition-opacity">
                    {brands.map((brand, i) => (
                        <span key={i} className="text-2xl font-black text-white tracking-tighter uppercase italic pointer-events-none select-none">
                            {brand}
                        </span>
                    ))}
                </div>
            </div>

            {/* Animated Gradient Line */}
            <div className="absolute bottom-0 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-blue-500/20 to-transparent"></div>
        </section>
    );
}

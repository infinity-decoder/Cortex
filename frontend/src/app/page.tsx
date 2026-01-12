'use client';

import { useState } from 'react';

export default function Home() {
  const [domain, setDomain] = useState('');
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<any>(null);

  const startScan = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/v1/scan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ domain }),
      });
      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error('Scan failed:', error);
      alert('Failed to run scan. Is the backend running?');
    }
    setLoading(false);
  };

  return (
    <main className="min-h-screen bg-slate-900 text-white p-8">
      <div className="max-w-4xl mx-auto">
        <header className="mb-12">
          <h1 className="text-4xl font-bold text-blue-400 mb-2">Cortex</h1>
          <p className="text-slate-400">Container-Aware Attack Surface Monitoring</p>
        </header>

        <section className="bg-slate-800 p-6 rounded-xl border border-slate-700 mb-8">
          <h2 className="text-xl font-semibold mb-4 text-slate-200">New Scan</h2>
          <div className="flex gap-4">
            <input
              type="text"
              placeholder="Enter domain (e.g., example.com)"
              className="flex-1 bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
            />
            <button
              onClick={startScan}
              disabled={loading || !domain}
              className="bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white font-medium py-2 px-6 rounded-lg transition-colors"
            >
              {loading ? 'Scanning...' : 'Start Scan'}
            </button>
          </div>
        </section>

        {results && (
          <div className="space-y-8 animate-in fade-in duration-500">
            <section>
              <h3 className="text-lg font-semibold mb-3 text-slate-300">Assets Discovered ({results.assets?.length || 0})</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                {results.assets?.map((asset: any, i: number) => (
                  <div key={i} className="bg-slate-800 p-3 rounded-lg border border-slate-700">
                    <span className="text-blue-400 font-mono">{asset.Subdomain || '@'}</span>
                    <span className="mx-2 text-slate-500">→</span>
                    <span className="text-slate-300 font-mono">{asset.IPs?.[0]}</span>
                  </div>
                ))}
              </div>
            </section>

            <section>
              <div className="flex justify-between items-end mb-3">
                <h3 className="text-lg font-semibold text-slate-300">Risk Findings ({results.all_findings?.length || 0})</h3>
                {results.new_findings?.length > 0 && (
                  <span className="text-xs bg-red-600/20 text-red-500 border border-red-600/50 px-2 py-0.5 rounded-full animate-pulse">
                    {results.new_findings.length} NEW FINDINGS
                  </span>
                )}
              </div>
              <div className="space-y-4">
                {results.all_findings?.map((finding: any, i: number) => {
                  const isNew = results.new_findings?.some((nf: any) => nf.type === finding.type && nf.severity === finding.severity);
                  return (
                    <div key={i} className={`p-4 rounded-lg border flex flex-col gap-2 relative transition-all hover:scale-[1.01] ${finding.severity === 'critical' ? 'bg-red-900/20 border-red-900/50' :
                        finding.severity === 'high' ? 'bg-orange-900/20 border-orange-900/50' :
                          'bg-slate-800 border-slate-700'
                      }`}>
                      {isNew && (
                        <div className="absolute -top-2 -right-2 bg-red-600 text-[10px] font-bold px-2 py-0.5 rounded shadow-lg animate-bounce">
                          NEW
                        </div>
                      )}
                      <div className="flex justify-between items-center">
                        <span className="font-bold text-lg uppercase tracking-tight">{finding.type}</span>
                        <span className={`text-xs font-bold px-2 py-1 rounded uppercase ${finding.severity === 'critical' ? 'bg-red-600' :
                            finding.severity === 'high' ? 'bg-orange-600' : 'bg-slate-600'
                          }`}>
                          {finding.severity}
                        </span>
                      </div>
                      <p className="text-slate-300 text-sm leading-relaxed">{finding.description}</p>
                      <div className="mt-2 p-2 bg-slate-900/50 rounded text-xs">
                        <span className="text-slate-500 font-bold mr-2">REMEDIATION:</span>
                        <span className="text-slate-300">{finding.remediation}</span>
                      </div>
                    </div>
                  );
                })}
                {(!results.all_findings || results.all_findings.length === 0) && (
                  <div className="text-slate-500 italic p-8 bg-slate-800/30 border border-dashed border-slate-700 rounded-xl text-center">
                    No critical exposures detected. Scan complete.
                  </div>
                )}
              </div>
            </section>
          </div>
        )}
      </div>
    </main>
  );
}

'use client';

import { useState, useEffect } from 'react';

export default function Home() {
  const [domain, setDomain] = useState('');
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<any>(null);
  const [stats, setStats] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [verifying, setVerifying] = useState(false);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/v1/stats');
      const data = await res.json();
      setStats(data);
    } catch (e) {
      console.error('Failed to fetch stats');
    }
  };

  const startScan = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch('http://localhost:8080/api/v1/scan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ domain }),
      });

      if (response.status === 403) {
        setError('Domain not verified. Ownership verification required.');
        setLoading(false);
        return;
      }

      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error('Scan failed:', error);
      setError('Failed to run scan. Is the backend running?');
    }
    setLoading(false);
  };

  const verifyDomain = async () => {
    setVerifying(true);
    try {
      const response = await fetch('http://localhost:8080/api/v1/domains/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ domain }),
      });
      if (response.ok) {
        alert('Domain verified successfully!');
        setError(null);
        startScan();
      } else {
        const data = await response.json();
        alert(`Verification failed. Please add TXT record: ${data.expected}`);
      }
    } catch (e) {
      alert('Verification error.');
    }
    setVerifying(false);
  };

  return (
    <main className="min-h-screen bg-slate-900 text-white p-8">
      <div className="max-w-4xl mx-auto">
        <header className="mb-12 flex justify-between items-start">
          <div>
            <h1 className="text-4xl font-bold text-blue-400 mb-2 tracking-tight">Cortex</h1>
            <p className="text-slate-400">Container-Aware Attack Surface Monitoring</p>
          </div>
          {stats && (
            <div className="grid grid-cols-2 gap-4 text-right">
              <div className="bg-slate-800/50 p-3 rounded-lg border border-slate-700">
                <div className="text-xs text-slate-500 uppercase font-bold">Total Assets</div>
                <div className="text-2xl font-bold text-slate-200">{stats.total_assets}</div>
              </div>
              <div className="bg-slate-800/50 p-3 rounded-lg border border-slate-700">
                <div className="text-xs text-slate-500 uppercase font-bold">Critical Risks</div>
                <div className="text-2xl font-bold text-red-500">{stats.critical_risks}</div>
              </div>
            </div>
          )}
        </header>

        <section className="bg-slate-800 p-6 rounded-xl border border-slate-700 mb-8 shadow-2xl">
          <h2 className="text-xl font-semibold mb-4 text-slate-200">New Scan</h2>
          <div className="flex gap-4">
            <input
              type="text"
              placeholder="Enter domain (e.g., example.com)"
              className="flex-1 bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
            />
            <button
              onClick={startScan}
              disabled={loading || !domain}
              className="bg-blue-600 hover:bg-blue-500 disabled:bg-slate-700 text-white font-medium py-2 px-6 rounded-lg transition-all shadow-lg active:scale-95"
            >
              {loading ? 'Scanning...' : 'Start Scan'}
            </button>
          </div>

          {error && (
            <div className="mt-4 p-4 bg-red-900/20 border border-red-900/50 rounded-lg flex justify-between items-center animate-in slide-in-from-top-2 duration-300">
              <span className="text-red-400 text-sm font-medium">{error}</span>
              {error.includes('verified') && (
                <button
                  onClick={verifyDomain}
                  disabled={verifying}
                  className="text-xs bg-red-600 hover:bg-red-500 px-3 py-1.5 rounded font-bold transition-colors"
                >
                  {verifying ? 'Verifying...' : 'Verify Now'}
                </button>
              )}
            </div>
          )}
        </section>

        {results && (
          <div className="space-y-8 animate-in fade-in duration-700 scale-in-95">
            <section>
              <h3 className="text-lg font-semibold mb-3 text-slate-300">Assets Discovered ({results.assets?.length || 0})</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                {results.assets?.map((asset: any, i: number) => (
                  <div key={i} className="bg-slate-800 p-3 rounded-lg border border-slate-700 flex items-center hover:bg-slate-700/50 transition-colors">
                    <div className="w-2 h-2 bg-green-500 rounded-full mr-3 animate-pulse"></div>
                    <span className="text-blue-400 font-mono text-sm">{asset.Subdomain || '@'}</span>
                    <span className="mx-2 text-slate-600">→</span>
                    <span className="text-slate-300 font-mono text-sm">{asset.IPs?.[0]}</span>
                  </div>
                ))}
              </div>
            </section>

            <section>
              <div className="flex justify-between items-end mb-3">
                <div className="flex items-center gap-3">
                  <h3 className="text-lg font-semibold text-slate-300">Risk Findings ({results.all_findings?.length || 0})</h3>
                  {results.verified && (
                    <span className="text-[10px] bg-green-900/40 text-green-500 border border-green-900/50 px-2 py-0.5 rounded tracking-widest font-bold uppercase">
                      Verified Scan
                    </span>
                  )}
                </div>
                {results.new_findings?.length > 0 && (
                  <span className="text-xs bg-red-600/20 text-red-500 border border-red-600/50 px-2 py-0.5 rounded-full animate-pulse font-bold">
                    {results.new_findings.length} NEW FINDINGS
                  </span>
                )}
              </div>
              <div className="space-y-4">
                {results.all_findings?.map((finding: any, i: number) => {
                  const isNew = results.new_findings?.some((nf: any) => nf.type === finding.type && nf.severity === finding.severity);
                  return (
                    <div key={i} className={`p-4 rounded-lg border flex flex-col gap-2 relative transition-all hover:scale-[1.01] ${finding.severity === 'critical' ? 'bg-red-900/30 border-red-900/50' :
                        finding.severity === 'high' ? 'bg-orange-900/20 border-orange-900/50' :
                          'bg-slate-800 border-slate-700'
                      }`}>
                      {isNew && (
                        <div className="absolute -top-2 -right-2 bg-red-600 text-[10px] font-bold px-2 py-0.5 rounded shadow-xl animate-bounce">
                          NEW
                        </div>
                      )}
                      <div className="flex justify-between items-center">
                        <span className="font-bold text-lg uppercase tracking-tight text-white">{finding.type}</span>
                        <span className={`text-[10px] font-bold px-2 py-1 rounded uppercase tracking-widest ${finding.severity === 'critical' ? 'bg-red-600' :
                            finding.severity === 'high' ? 'bg-orange-600' : 'bg-slate-600'
                          }`}>
                          {finding.severity}
                        </span>
                      </div>
                      <p className="text-slate-300 text-sm leading-relaxed">{finding.description}</p>
                      <div className="mt-2 p-2 bg-slate-900/50 rounded text-xs border border-white/5">
                        <span className="text-slate-500 font-bold mr-2">REMEDIATION:</span>
                        <span className="text-slate-300 italic">{finding.remediation}</span>
                      </div>
                    </div>
                  );
                })}
                {(!results.all_findings || results.all_findings.length === 0) && (
                  <div className="text-slate-500 italic p-12 bg-slate-800/20 border border-dashed border-slate-700 rounded-2xl text-center">
                    No critical exposures detected. Your attack surface appears clean.
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

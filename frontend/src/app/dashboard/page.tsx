'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import DashboardLayout from '@/components/DashboardLayout';
import {
  ShieldAlert,
  Globe,
  Zap,
  Activity,
  ArrowUpRight,
  Plus,
  Play
} from 'lucide-react';

export default function Dashboard() {
  const [domain, setDomain] = useState('');
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [assets, setAssets] = useState<any[]>([]);
  const [findings, setFindings] = useState<any[]>([]);
  const [stats, setStats] = useState<any>(null);
  const [verifiedDomains, setVerifiedDomains] = useState<any[]>([]);

  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }
    fetchDomains();
  }, []);

  // Effect to load data once a domain is selected
  useEffect(() => {
    if (selectedDomain) {
      fetchStats();
      fetchAssets();
      fetchFindings();
    }
  }, [selectedDomain]);

  const getHeaders = () => {
    const token = localStorage.getItem('token');
    return {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    };
  };

  const fetchDomains = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/v1/domains', { headers: getHeaders() });
      if (res.ok) {
        const data = await res.json();
        setVerifiedDomains(data || []);
        // Auto-select first domain if available
        if (data && data.length > 0) {
          setSelectedDomain(data[0].RootDomain);
        }
      }
    } catch (e) {
      console.warn('Failed to fetch domains');
    }
  };

  const fetchStats = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/v1/stats', {
        headers: getHeaders()
      });
      if (!res.ok) throw new Error('API Error');
      const data = await res.json();
      setStats(data);
    } catch (e) {
      setStats(null);
    }
  };

  const fetchAssets = async () => {
    if (!selectedDomain) return;
    try {
      const res = await fetch(`http://localhost:8080/api/v1/assets?domain=${selectedDomain}`, {
        headers: getHeaders()
      });
      if (!res.ok) throw new Error('API Error');
      const data = await res.json();
      setAssets(data.slice(0, 4));
    } catch (e) {
      setAssets([]);
    }
  };

  const fetchFindings = async () => {
    if (!selectedDomain) return;
    try {
      const res = await fetch(`http://localhost:8080/api/v1/findings?domain=${selectedDomain}`, {
        headers: getHeaders()
      });
      if (!res.ok) throw new Error('API Error');
      const data = await res.json();
      setFindings(data.slice(0, 3));
    } catch (e) {
      setFindings([]);
    }
  };

  return (
    <DashboardLayout>
      <div className="max-w-7xl mx-auto space-y-10">
        {/* Header Section */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-extrabold text-white tracking-tight mb-1">Security Overview</h1>
            <p className="text-slate-500 text-sm font-medium uppercase tracking-widest">
              {selectedDomain ? `Monitoring: ${selectedDomain}` : 'No Domains Configured'}
            </p>
          </div>
          <div className="flex items-center gap-3">
            {/* Domain Selector */}
            {verifiedDomains.length > 0 && (
              <select
                className="bg-slate-900 border border-slate-800 text-white text-xs font-bold rounded-xl px-4 py-3 outline-none focus:ring-1 focus:ring-blue-500"
                value={selectedDomain || ''}
                onChange={(e) => setSelectedDomain(e.target.value)}
              >
                {verifiedDomains.map((d: any) => (
                  <option key={d.ID} value={d.RootDomain}>{d.RootDomain}</option>
                ))}
              </select>
            )}

            <div className="bg-slate-900 border border-slate-800 rounded-xl p-1 flex gap-1">
              <button className="px-4 py-2 bg-slate-800 text-white rounded-lg text-xs font-bold shadow-sm">24H</button>
            </div>
            <button className="bg-blue-600 hover:bg-blue-500 text-white p-3 rounded-xl transition-all shadow-lg shadow-blue-500/10 active:scale-95">
              <Plus size={20} />
            </button>
          </div>
        </div>

        {/* Empty State if No Domain */}
        {!selectedDomain ? (
          <div className="bg-slate-900/40 border border-slate-800 rounded-2xl p-12 text-center">
            <Globe size={48} className="mx-auto text-slate-600 mb-6" />
            <h3 className="text-xl font-bold text-white mb-2">No Verified Domains</h3>
            <p className="text-slate-500 max-w-md mx-auto mb-8">
              You haven't verified any domains yet. Add a domain to start continuous monitoring.
            </p>
            <div className="flex gap-4 justify-center">
              <input
                type="text"
                placeholder="example.com"
                className="bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-white focus:outline-none focus:ring-1 focus:ring-blue-500"
                value={domain}
                onChange={(e) => setDomain(e.target.value)}
              />
              <button className="bg-blue-600 text-white px-6 py-3 rounded-xl font-bold hover:bg-blue-500 transition-all">
                Add Domain
              </button>
            </div>
          </div>
        ) : (
          <>
            {/* Scorecards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl relative overflow-hidden group">
                <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                  <Globe size={64} />
                </div>
                <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">Total Assets</div>
                <div className="text-3xl font-extrabold text-white mb-2">{stats?.total_assets || 0}</div>
                <div className="flex items-center gap-1.5 text-green-500 text-xs font-bold">
                  <ArrowUpRight size={14} />
                  <span>Live</span>
                </div>
              </div>

              <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl relative overflow-hidden group">
                <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                  <Zap size={64} />
                </div>
                <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">Exposed Services</div>
                <div className="text-3xl font-extrabold text-white mb-2">{stats?.scans_completed || 0}</div>
                <div className="flex items-center gap-1.5 text-orange-500 text-xs font-bold">
                  <Activity size={14} className="animate-pulse" />
                  <span>Active Endpoints</span>
                </div>
              </div>

              <div className="bg-red-500/5 border border-red-500/20 p-6 rounded-2xl relative overflow-hidden group">
                <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                  <ShieldAlert size={64} className="text-red-500" />
                </div>
                <div className="text-[10px] font-extrabold text-red-500/50 uppercase tracking-widest mb-1">Critical Risks</div>
                <div className="text-3xl font-extrabold text-red-500 mb-2">{stats?.critical_risks || 0}</div>
                <div className="flex items-center gap-1.5 text-red-400 text-xs font-bold">
                  <span>Attention Required</span>
                </div>
              </div>

              <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-2xl relative overflow-hidden group">
                <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">System Risk Score</div>
                <div className="text-3xl font-extrabold text-white mb-2">84<span className="text-slate-600 text-xl">/100</span></div>
                <div className="flex items-center gap-1.5 text-slate-500 text-xs font-bold">
                  <span>Industry Avg: 65</span>
                </div>
              </div>
            </div>

            {/* Main Sections */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              <div className="lg:col-span-2 space-y-8">
                {/* Quick Scan Section */}
                <section className="bg-slate-900/40 border border-slate-800 rounded-2xl p-8 relative overflow-hidden">
                  <div className="relative z-10">
                    <h2 className="text-xl font-bold text-white mb-4 flex items-center gap-3">
                      <Play size={20} className="text-blue-500 fill-current" />
                      Trigger Manual Recon
                    </h2>
                    <div className="flex gap-4">
                      <input
                        type="text"
                        placeholder="Enter domain (e.g., prod-api.cortex.security)"
                        className="flex-1 bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 text-sm font-medium text-white placeholder:text-slate-600 focus:outline-none focus:ring-1 focus:ring-blue-500/50 transition-all"
                        value={domain}
                        onChange={(e) => setDomain(e.target.value)}
                      />
                      <button
                        disabled={loading || !domain}
                        className="bg-blue-600 hover:bg-blue-500 disabled:bg-slate-800 text-white font-bold py-3 px-8 rounded-xl transition-all shadow-lg active:scale-95 flex items-center gap-2"
                      >
                        {loading ? 'Analyzing...' : 'Execute Scan'}
                      </button>
                    </div>
                  </div>
                </section>

                {/* Subdomain Explorer */}
                <section>
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-bold text-white flex items-center gap-3">
                      <Globe size={20} className="text-slate-400" />
                      Active Monitoring Subdomains
                    </h3>
                    <button className="text-xs text-blue-500 font-bold hover:text-blue-400">View All Assets</button>
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {assets.length > 0 ? assets.map((asset, i) => (
                      <div key={i} className="bg-slate-900/50 border border-slate-800 p-4 rounded-xl flex items-center justify-between group hover:border-slate-700 transition-colors">
                        <div className="flex items-center gap-3">
                          <div className="w-2 h-2 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]"></div>
                          <div>
                            <div className="text-sm font-bold text-white font-mono">{asset.Subdomain || '@'}</div>
                            <div className="text-[10px] text-slate-500 font-mono mt-0.5">{asset.IpAddress || 'Unresolved'}</div>
                          </div>
                        </div>
                        <ArrowUpRight size={14} className="text-slate-600 group-hover:text-white transition-colors" />
                      </div>
                    )) : (
                      <div className="col-span-2 py-10 text-center bg-slate-900/20 border border-slate-800 rounded-xl border-dashed">
                        <p className="text-xs font-bold text-slate-600 uppercase tracking-widest">No active monitoring data</p>
                      </div>
                    )}
                  </div>
                </section>
              </div>

              {/* Sidebar Widgets */}
              <div className="space-y-8">
                <section className="bg-slate-900/40 border border-slate-800 rounded-2xl p-6">
                  <h3 className="text-sm font-extrabold text-slate-500 uppercase tracking-widest mb-6">Recent Delta Alerts</h3>
                  <div className="space-y-6">
                    {findings.length > 0 ? findings.map((finding, i) => (
                      <div key={i} className="flex gap-4 relative">
                        <div className="flex flex-col items-center">
                          <div className={`w-3 h-3 rounded-full ${(finding.Severity || finding.severity) === 'critical' ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.5)]' :
                            (finding.Severity || finding.severity) === 'high' ? 'bg-orange-500' : 'bg-blue-500'
                            }`}></div>
                          {i < findings.length - 1 && <div className="w-px flex-1 bg-slate-800 my-2"></div>}
                        </div>
                        <div className="-mt-1 pb-4">
                          <div className="text-sm font-bold text-white leading-tight">{finding.Type || finding.title}</div>
                          <div className="text-[10px] text-slate-600 font-medium mt-1 uppercase tracking-wider">{finding.Detected || 'Just now'}</div>
                        </div>
                      </div>
                    )) : (
                      <div className="text-center py-4">
                        <p className="text-[10px] font-bold text-slate-600 uppercase tracking-widest">All systems nominal</p>
                      </div>
                    )}
                  </div>
                </section>

                <section className="bg-blue-600 p-6 rounded-2xl shadow-xl shadow-blue-600/10">
                  <h3 className="text-white font-bold mb-2">Upgrade to Pro</h3>
                  <p className="text-blue-100 text-xs leading-relaxed mb-6">
                    Get hourly scans, team collaboration, and real-time Slack integrations.
                  </p>
                  <button className="w-full bg-slate-900 text-white text-xs font-extrabold py-3 rounded-xl hover:bg-slate-800 transition-all shadow-lg active:scale-95 uppercase tracking-widest">
                    Learn More
                  </button>
                </section>
              </div>
            </div>
          </>
        )}
      </div>
    </DashboardLayout>
  );
}

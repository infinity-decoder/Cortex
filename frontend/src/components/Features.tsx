import { Server, Layout, Database, Workflow, CheckCircle, Globe } from 'lucide-react';

const monitoringItems = [
    {
        title: 'Kubernetes Nodes',
        description: 'Detection of exposed Kubelets and API servers across clusters.',
        icon: Layout,
    },
    {
        title: 'Docker Remote API',
        description: 'Identify unauthenticated TCP or TLS-misconfigured Docker endpoints.',
        icon: Server,
    },
    {
        title: 'Container Registries',
        description: 'Monitoring for public /v2/_catalog exposure and leaky images.',
        icon: Database,
    },
    {
        title: 'Ingress Controllers',
        description: 'Fingerprinting NGINX, Istio, and Traefik dashboard leaks.',
        icon: Workflow,
    },
    {
        title: 'Custom Root Domains',
        description: 'Recursive DNS enumeration and TLS certificate analysis.',
        icon: Globe,
    },
    {
        title: 'Continuous Verification',
        description: 'Automated ownership checks and scope enforcement.',
        icon: CheckCircle,
    },
];

export default function Features() {
    return (
        <section id="features" className="py-24 bg-slate-900/50 relative">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="text-center mb-16">
                    <h2 className="text-3xl md:text-5xl font-bold text-white mb-4">What Cortex Monitors</h2>
                    <p className="text-slate-400 max-w-2xl mx-auto">
                        Unlike generic port scanners, Cortex is purpose-built for the modern cloud-native stack. We understand container orchestration and internal service leakage.
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    {monitoringItems.map((item, i) => (
                        <div key={i} className="p-8 rounded-2xl bg-slate-800/40 border border-slate-700/50 hover:bg-slate-800/60 transition-all group">
                            <div className="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center mb-6 text-blue-500 group-hover:bg-blue-500 group-hover:text-white transition-all">
                                <item.icon size={24} />
                            </div>
                            <h3 className="text-xl font-bold text-white mb-3">{item.title}</h3>
                            <p className="text-slate-400 leading-relaxed text-sm">
                                {item.description}
                            </p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}

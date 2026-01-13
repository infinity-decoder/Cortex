'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
    BarChart3,
    ShieldAlert,
    Globe,
    Clock,
    Settings,
    ChevronRight,
    LogOut,
    Zap
} from 'lucide-react';

const menuItems = [
    { name: 'Overview', icon: BarChart3, href: '/dashboard' },
    { name: 'Assets', icon: Globe, href: '/dashboard/assets' },
    { name: 'Findings', icon: ShieldAlert, href: '/dashboard/findings' },
    { name: 'Monitoring', icon: Clock, href: '/dashboard/monitoring' },
    { name: 'Settings', icon: Settings, href: '/dashboard/settings' },
];

export default function Sidebar() {
    const pathname = usePathname();

    return (
        <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col h-screen sticky top-0">
            <div className="p-6">
                <Link href="/dashboard" className="flex items-center gap-2 mb-8">
                    <div className="w-8 h-8 bg-blue-600 rounded flex items-center justify-center">
                        <Zap size={18} className="text-white fill-current" />
                    </div>
                    <span className="text-xl font-bold tracking-tight text-white">Cortex</span>
                </Link>

                <nav className="space-y-1">
                    {menuItems.map((item) => {
                        const isActive = pathname === item.href;
                        return (
                            <Link
                                key={item.name}
                                href={item.href}
                                className={`flex items-center justify-between px-4 py-3 rounded-xl transition-all group ${isActive
                                        ? 'bg-blue-600/10 text-blue-400 border border-blue-500/20 shadow-lg shadow-blue-500/5'
                                        : 'text-slate-500 hover:text-slate-300 hover:bg-slate-800'
                                    }`}
                            >
                                <div className="flex items-center gap-3">
                                    <item.icon size={20} className={isActive ? 'text-blue-400' : 'group-hover:text-slate-300'} />
                                    <span className="text-sm font-bold">{item.name}</span>
                                </div>
                                {isActive && <ChevronRight size={14} />}
                            </Link>
                        );
                    })}
                </nav>
            </div>

            <div className="mt-auto p-6 border-t border-slate-800">
                <div className="bg-slate-800/40 p-4 rounded-xl border border-slate-700/50 mb-4">
                    <div className="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest mb-1">Current Plan</div>
                    <div className="text-sm font-bold text-white mb-2 uppercase tracking-tight">Free Tier</div>
                    <div className="w-full bg-slate-700 h-1 rounded-full overflow-hidden">
                        <div className="bg-blue-500 h-full w-1/5"></div>
                    </div>
                    <div className="text-[10px] text-slate-500 mt-2 font-medium">1 / 5 Domains Used</div>
                </div>

                <button className="flex items-center gap-3 px-4 py-3 w-full text-slate-500 hover:text-red-400 hover:bg-red-500/10 rounded-xl transition-all group">
                    <LogOut size={20} className="group-hover:text-red-400" />
                    <span className="text-sm font-bold">Logout</span>
                </button>
            </div>
        </aside>
    );
}

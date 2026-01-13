<?php
/**
 * Cortex Smart Gateway
 * Intelligent routing between Installer, Landing Page, and Dashboard.
 */

// 1. Bootstrap Check
$is_installed = file_exists('.env');
if (!$is_installed) {
    header('Location: ./install/');
    exit;
}

// 2. Routing Configuration
$frontend_base = "http://localhost:3000";
$is_logged_in = isset($_COOKIE['token']) || isset($_COOKIE['auth_token']);
$target_url = $is_logged_in ? "$frontend_base/dashboard" : "$frontend_base/";

/**
 * We render a premium transition page to ensure the user understands 
 * the bridge between the PHP environment and the Next.js Frontend.
 */
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cortex | Smart Gateway</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700;800&display=swap" rel="stylesheet">
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        body { font-family: 'Inter', sans-serif; }
        .bg-mesh {
            background-color: #0f172a;
            background-image: 
                radial-gradient(at 0% 0%, rgba(30, 64, 175, 0.15) 0px, transparent 50%),
                radial-gradient(at 100% 100%, rgba(30, 58, 138, 0.15) 0px, transparent 50%);
        }
    </style>
    <meta http-equiv="refresh" content="1;url=<?php echo $target_url; ?>">
</head>
<body class="bg-mesh text-slate-200 min-h-screen flex items-center justify-center p-6">
    <div class="max-w-md w-full bg-slate-900/50 border border-slate-800 rounded-3xl p-10 text-center shadow-2xl backdrop-blur-xl">
        <div class="w-20 h-20 bg-blue-600 rounded-2xl mx-auto mb-8 flex items-center justify-center shadow-lg shadow-blue-500/20">
            <svg xmlns="http://www.w3.org/2000/svg" class="text-white w-12 h-12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
        </div>
        
        <h1 class="text-3xl font-extrabold text-white mb-3 tracking-tight italic uppercase">Cortex Gateway</h1>
        <p class="text-slate-500 text-sm mb-6 leading-relaxed font-medium">
            Bridging to your high-performance security surface. 
        </p>

        <div class="bg-blue-600/10 border border-blue-500/20 p-4 rounded-2xl mb-10 text-left">
            <h4 class="text-[10px] font-black text-blue-500 uppercase tracking-widest mb-2">Network Topology</h4>
            <div class="flex items-center justify-between text-[11px] font-mono">
                <span class="text-slate-400">Apache Gateway</span>
                <span class="text-blue-500 font-bold">Port 80</span>
            </div>
            <div class="my-2 border-t border-blue-500/10"></div>
            <div class="flex items-center justify-between text-[11px] font-mono">
                <span class="text-slate-400">Next.js Core</span>
                <span class="text-green-500 font-bold">Port 3000</span>
            </div>
        </div>

        <div class="space-y-4">
            <p class="text-xs text-slate-400 italic">Redirecting to <?php echo $is_logged_in ? 'Dashboard' : 'Landing Page'; ?>...</p>
            <a href="<?php echo $target_url; ?>" 
               class="block w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-4 rounded-2xl transition-all shadow-lg shadow-blue-500/10 active:scale-95 text-sm uppercase tracking-widest">
                Continue to Platform
            </a>
            
            <div class="flex items-center justify-center gap-2">
                <div class="w-1 h-1 rounded-full bg-blue-500 animate-pulse"></div>
                <div class="w-1 h-1 rounded-full bg-blue-500 animate-pulse delay-75"></div>
                <div class="w-1 h-1 rounded-full bg-blue-500 animate-pulse delay-150"></div>
            </div>
        </div>

        <div class="mt-12 pt-8 border-t border-slate-800/50">
            <p class="text-[10px] font-bold text-slate-600 uppercase tracking-widest">
                Environment: Development (Localhost:3000)
            </p>
        </div>
    </div>
</body>
</html>

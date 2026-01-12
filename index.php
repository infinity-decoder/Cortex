<?php
/**
 * Cortex Entry Point
 * Handles redirection between Installer and Landing Page
 */

// Check for existence of .env file
$is_installed = file_exists('.env');

if (!$is_installed) {
    // Not installed, redirect to the Setup Wizard
    header('Location: ./install/');
    exit;
}

// Already installed. 
// In a production environment with a proxy, this part would usually be handled by Nginx/Apache.
// For local testing/Laragon, we'll redirect to the Frontend Dashboard.
// We assume default port 3000 for Next.js, or the user can navigate manually.

?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cortex | Security Monitoring</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" rel="stylesheet">
    <style>
        body {
            background: #0f172a;
            color: white;
            font-family: 'Inter', sans-serif;
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
            text-align: center;
        }
        .container {
            max-width: 600px;
            padding: 2rem;
            background: #1e293b;
            border-radius: 1rem;
            border: 1px solid #334155;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
        }
        img {
            width: 80px;
            margin-bottom: 1.5rem;
        }
        h1 {
            font-size: 2rem;
            margin-bottom: 1rem;
        }
        p {
            color: #94a3b8;
            line-height: 1.6;
            margin-bottom: 2rem;
        }
        .btn-group {
            display: flex;
            gap: 1rem;
            justify-content: center;
        }
        .btn {
            background: #3b82f6;
            color: white;
            text-decoration: none;
            padding: 0.75rem 1.5rem;
            border-radius: 0.5rem;
            font-weight: 600;
            transition: 0.2s;
        }
        .btn:hover {
            filter: brightness(1.1);
        }
        .btn-secondary {
            background: #334155;
        }
    </style>
</head>
<body>
    <div class="container">
        <img src="frontend/public/cortex_SaaS.png" alt="Cortex Logo">
        <h1>Welcome to Cortex</h1>
        <p>Your attack surface monitoring platform is configured and ready. Please ensure your Backend and Frontend services are running to access the dashboard.</p>
        <div class="btn-group">
            <a href="http://localhost:3000" class="btn">Launch Dashboard</a>
            <a href="https://github.com/infinity-decoder/Cortex" target="_blank" class="btn btn-secondary">View Repository</a>
        </div>
    </div>
</body>
</html>

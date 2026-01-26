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

// 2. Serve the Application
// We serve the statically exported Next.js app (app.html)
if (file_exists('app.html')) {
    readfile('app.html');
    exit;
} else {
    // Fallback if build is missing
    echo "Cortex Frontend build not found. Please run 'npm run build' in the frontend directory and copy 'out' contents to root.";
    exit;
}
?>

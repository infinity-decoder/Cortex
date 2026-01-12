<?php
/**
 * Cortex Setup Wizard
 * Inspired by Shiori Installer
 */

$is_installed = file_exists('../.env');
if ($is_installed && !isset($_GET['reinstall'])) {
    header('Location: /');
    exit;
}

?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cortex | Setup Wizard</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #0f172a;
            --accent: #3b82f6;
            --card: #1e293b;
            --text: #f8fafc;
        }
        body {
            background: var(--bg);
            color: var(--text);
            font-family: 'Inter', sans-serif;
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
        }
        .container {
            background: var(--card);
            padding: 2.5rem;
            border-radius: 1rem;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
            width: 100%;
            max-width: 500px;
            border: 1px solid rgba(255,255,255,0.05);
        }
        .header {
            text-align: center;
            margin-bottom: 2rem;
        }
        .header img {
            width: 64px;
            margin-bottom: 1rem;
        }
        h1 {
            font-size: 1.5rem;
            font-weight: 800;
            margin: 0;
            letter-spacing: -0.025em;
        }
        .step {
            display: none;
        }
        .step.active {
            display: block;
        }
        .input-group {
            margin-bottom: 1.25rem;
        }
        label {
            display: block;
            font-size: 0.875rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
            color: #94a3b8;
        }
        input {
            width: 100%;
            background: var(--bg);
            border: 1px solid #334155;
            color: white;
            padding: 0.75rem;
            border-radius: 0.5rem;
            font-size: 0.875rem;
            box-sizing: border-box;
        }
        input:focus {
            outline: none;
            border-color: var(--accent);
            ring: 2px solid var(--accent);
        }
        .btn {
            background: var(--accent);
            color: white;
            border: none;
            width: 100%;
            padding: 0.875rem;
            border-radius: 0.5rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
            margin-top: 1rem;
        }
        .btn:hover {
            filter: brightness(1.1);
            transform: translateY(-1px);
        }
        .secondary-btn {
            background: transparent;
            color: #94a3b8;
            margin-top: 0.5rem;
            font-size: 0.75rem;
            border: 1px solid #334155;
        }
        .loader {
            display: none;
            border: 3px solid #f3f3f3;
            border-top: 3px solid var(--accent);
            border-radius: 50%;
            width: 20px;
            height: 20px;
            animation: spin 1s linear infinite;
            margin: 10px auto;
        }
        @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="https://raw.githubusercontent.com/infinity-decoder/Cortex/master/frontend/public/cortex_SaaS.png" alt="Logo">
            <h1>Cortex Setup Wizard</h1>
        </div>

        <form id="setupForm" action="process.php" method="POST">
            <!-- Step 1: Database -->
            <div id="step1" class="step active">
                <h3>Database Configuration</h3>
                <div class="input-group">
                    <label>PostgreSQL Host</label>
                    <input type="text" name="db_host" value="localhost" required>
                </div>
                <div class="input-group">
                    <label>PostgreSQL Port</label>
                    <input type="text" name="db_port" value="5432" required>
                </div>
                <div class="input-group">
                    <label>Database User</label>
                    <input type="text" name="db_user" placeholder="postgres" required>
                </div>
                <div class="input-group">
                    <label>Database Password</label>
                    <input type="password" name="db_pass" placeholder="********">
                </div>
                <div class="input-group">
                    <label>Database Name</label>
                    <input type="text" name="db_name" value="cortex" required>
                </div>
                <button type="button" class="btn" onclick="nextStep(2)">Configure Application</button>
            </div>

            <!-- Step 2: App Config -->
            <div id="step2" class="step">
                <h3>Application Settings</h3>
                <div class="input-group">
                    <label>Admin Email</label>
                    <input type="email" name="admin_email" placeholder="admin@example.com" required>
                </div>
                <div class="input-group">
                    <label>Organization Name</label>
                    <input type="text" name="org_name" value="Default Org" required>
                </div>
                <button type="submit" class="btn">Finish Installation</button>
                <button type="button" class="btn secondary-btn" onclick="nextStep(1)">Back</button>
            </div>
            
            <div id="loading" class="loader"></div>
        </form>
    </div>

    <script>
        function nextStep(step) {
            document.querySelectorAll('.step').forEach(s => s.classList.remove('active'));
            document.getElementById('step' + step).classList.add('active');
        }

        document.getElementById('setupForm').onsubmit = function() {
            document.getElementById('loading').style.display = 'block';
            document.querySelector('.btn').disabled = true;
        };
    </script>
</body>
</html>

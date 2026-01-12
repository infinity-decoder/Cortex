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
    <script src="https://cdn.jsdelivr.net/npm/sweetalert2@11"></script>
    <style>
        :root {
            --bg: #0f172a;
            --accent: #3b82f6;
            --card: #1e293b;
            --text: #f8fafc;
            --error: #ef4444;
            --success: #10b981;
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
            transition: transform 0.3s ease;
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
            animation: fadeIn 0.4s ease;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
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
            transition: all 0.2s;
        }
        input:focus {
            outline: none;
            border-color: var(--accent);
            box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
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
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
        }
        .btn:hover:not(:disabled) {
            filter: brightness(1.1);
            transform: translateY(-1px);
        }
        .btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
        .secondary-btn {
            background: transparent;
            color: #94a3b8;
            margin-top: 0.5rem;
            font-size: 0.75rem;
            border: 1px solid #334155;
        }
        .swal2-popup {
            background: var(--card) !important;
            color: var(--text) !important;
            border-radius: 1rem !important;
            border: 1px solid rgba(255,255,255,0.1) !important;
        }
        .hint {
            font-size: 0.75rem;
            color: #64748b;
            margin-top: 0.25rem;
        }
    </style>
</head>
<body>
    <div class="container" id="installerContainer">
        <div class="header">
            <img src="https://raw.githubusercontent.com/infinity-decoder/Cortex/master/frontend/public/cortex_SaaS.png" alt="Logo">
            <h1>Cortex Setup Wizard</h1>
        </div>

        <form id="setupForm">
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
                    <input type="text" name="db_user" placeholder="postgres" value="postgres" required>
                    <div class="hint">Recommended: Use 'postgres' or a user with CREATEDB privileges.</div>
                </div>
                <div class="input-group">
                    <label>Database Password</label>
                    <input type="password" name="db_pass" placeholder="********">
                    <div class="hint">Laragon default is often empty.</div>
                </div>
                <div class="input-group">
                    <label>Database Name</label>
                    <input type="text" name="db_name" value="cortex" required>
                    <div class="hint">Installer will attempt to create this database.</div>
                </div>
                <button type="button" class="btn" onclick="nextStep(2)">Configure Account &rarr;</button>
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
                    <input type="text" name="org_name" value="My Organization" required>
                </div>
                <button type="submit" class="btn" id="submitBtn">Install Cortex</button>
                <button type="button" class="btn secondary-btn" onclick="nextStep(1)">Back</button>
            </div>
        </form>
    </div>

    <script>
        function nextStep(step) {
            document.querySelectorAll('.step').forEach(s => s.classList.remove('active'));
            document.getElementById('step' + step).classList.add('active');
        }

        document.getElementById('setupForm').onsubmit = async function(e) {
            e.preventDefault();
            
            const btn = document.getElementById('submitBtn');
            const originalText = btn.innerHTML;
            btn.disabled = true;
            btn.innerHTML = 'Installing...';

            const formData = new FormData(this);

            try {
                const response = await fetch('process.php', {
                    method: 'POST',
                    body: formData
                });

                const result = await response.json();

                if (result.success) {
                    Swal.fire({
                        icon: 'success',
                        title: 'Installation Successful!',
                        text: result.message,
                        confirmButtonText: 'Launch Portal',
                        confirmButtonColor: '#3b82f6'
                    }).then(() => {
                        window.location.href = '/';
                    });
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Installation Failed',
                        text: result.message,
                        confirmButtonColor: '#ef4444'
                    });
                }
            } catch (error) {
                Swal.fire({
                    icon: 'error',
                    title: 'System Error',
                    text: 'An unexpected error occurred during installation.',
                    footer: error.message
                });
            } finally {
                btn.disabled = false;
                btn.innerHTML = originalText;
            }
        };
    </script>
</body>
</html>

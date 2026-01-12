<?php
/**
 * Cortex Installation Processor
 */

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    header('Location: index.php');
    exit;
}

// 0. Sanitize Inputs
$db_host = filter_input(INPUT_POST, 'db_host', FILTER_SANITIZE_FULL_SPECIAL_CHARS);
$db_port = filter_input(INPUT_POST, 'db_port', FILTER_SANITIZE_NUMBER_INT);
$db_user = filter_input(INPUT_POST, 'db_user', FILTER_SANITIZE_FULL_SPECIAL_CHARS);
$db_pass = $_POST['db_pass']; // Passwords shouldn't be HTML sanitized
$db_name = filter_input(INPUT_POST, 'db_name', FILTER_SANITIZE_FULL_SPECIAL_CHARS);
$admin_email = filter_input(INPUT_POST, 'admin_email', FILTER_SANITIZE_EMAIL);
$org_name = filter_input(INPUT_POST, 'org_name', FILTER_SANITIZE_FULL_SPECIAL_CHARS);

// Validation
if (!$admin_email || !filter_var($admin_email, FILTER_VALIDATE_EMAIL)) {
    die("Invalid admin email address.");
}

// 1. Test PDO Connection (Corrected constant from ATTR_ERR_CODE to ATTR_ERRMODE)
try {
    $dsn = "pgsql:host=$db_host;port=$db_port;dbname=postgres";
    $pdo = new PDO($dsn, $db_user, $db_pass, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);

    // Create DB if not exists
    $pdo->exec("CREATE DATABASE $db_name");
} catch (PDOException $e) {
    if (strpos($e->getMessage(), 'already exists') === false) {
        die("Database connection failed: " . $e->getMessage());
    }
}

// 2. Run Migrations
try {
    $dsn = "pgsql:host=$db_host;port=$db_port;dbname=$db_name";
    $pdo = new PDO($dsn, $db_user, $db_pass, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);

    $sql = file_get_contents('../backend/schema.sql');
    $pdo->exec($sql);

    // Initial Org & Domain
    $org_id = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(16384, 20479), mt_rand(32768, 49151), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535));
    $pdo->prepare("INSERT INTO organizations (id, name, plan_tier) VALUES (?, ?, 'free')")->execute([$org_id, $org_name]);
    
    // Admin User Placeholder (Password: admin123)
    $user_id = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(16384, 20479), mt_rand(32768, 49151), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535));
    $pdo->prepare("INSERT INTO users (id, email, password_hash) VALUES (?, ?, '$2a$10$7E1gPFF1z6Y9f6Y9f6Y9f6Y9f6Y9f6Y9f6Y9f6Y9f6Y9f6')")->execute([$user_id, $admin_email]);

} catch (PDOException $e) {
    die("Migration failed: " . $e->getMessage());
}

// 3. Write .env File
$env_content = "DATABASE_URL=postgres://$db_user:$db_pass@$db_host:$db_port/$db_name?sslmode=disable\n";
$env_content .= "PORT=8080\n";
$env_content .= "NODE_ENV=production\n";

file_put_contents('../.env', $env_content);

// 4. Success Page
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Cortex | Installation Successful</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" rel="stylesheet">
    <style>
        body { background: #0f172a; color: white; font-family: 'Inter', sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; text-align: center; }
        .card { background: #1e293b; padding: 3rem; border-radius: 1rem; border: 1px solid #334155; max-width: 400px; }
        h2 { color: #10b981; }
        p { color: #94a3b8; line-height: 1.6; }
        .btn { display: inline-block; background: #3b82f6; color: white; text-decoration: none; padding: 0.75rem 1.5rem; border-radius: 0.5rem; font-weight: 600; margin-top: 1.5rem; }
    </style>
</head>
<body>
    <div class="card">
        <h2>Installation Successful!</h2>
        <p>Cortex has been configured and is ready to secure your attack surface.</p>
        <p><strong>Next Steps:</strong><br/>1. Start the Go Backend (bin/api.exe)<br/>2. Start the Frontend (npm run dev)</p>
        <a href="/" class="btn">Launch Portal</a>
    </div>
</body>
</html>

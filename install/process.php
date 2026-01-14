<?php
/**
 * Cortex Installation Processor (AJAX Version)
 */

header('Content-Type: application/json');

function sendResponse($success, $message, $data = []) {
    echo json_encode(['success' => $success, 'message' => $message, 'data' => $data]);
    exit;
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    sendResponse(false, 'Invalid request method.');
}

// 0. Sanitize Inputs
$db_host = filter_input(INPUT_POST, 'db_host', FILTER_SANITIZE_FULL_SPECIAL_CHARS);
$db_port = filter_input(INPUT_POST, 'db_port', FILTER_SANITIZE_NUMBER_INT) ?: 5432;
$db_user = filter_input(INPUT_POST, 'db_user', FILTER_SANITIZE_FULL_SPECIAL_CHARS);
$db_pass = $_POST['db_pass']; 
$db_name = filter_input(INPUT_POST, 'db_name', FILTER_SANITIZE_FULL_SPECIAL_CHARS);
$admin_email = filter_input(INPUT_POST, 'admin_email', FILTER_SANITIZE_EMAIL);
$org_name = filter_input(INPUT_POST, 'org_name', FILTER_SANITIZE_FULL_SPECIAL_CHARS);

// Validate database name (alphanumeric and underscore only)
if (!preg_match('/^[a-zA-Z0-9_]+$/', $db_name)) {
    sendResponse(false, 'Invalid database name. Only alphanumeric characters and underscores allowed.');
}

// Validate email
if (!filter_var($admin_email, FILTER_VALIDATE_EMAIL)) {
    sendResponse(false, 'Invalid email address.');
}

// 1. Test Connection & Create Database
try {
    // Connect to 'postgres' system database first to create the target database
    $dsn = "pgsql:host=$db_host;port=$db_port;dbname=postgres";
    $pdo = new PDO($dsn, $db_user, $db_pass, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);

    // Check if database exists
    $stmt = $pdo->prepare("SELECT 1 FROM pg_database WHERE datname = ?");
    $stmt->execute([$db_name]);
    if (!$stmt->fetch()) {
        // Escape database name using pg_escape_identifier equivalent
        // Since PDO doesn't have pg_escape_identifier, we validate and quote manually
        $db_name_quoted = '"' . str_replace('"', '""', $db_name) . '"';
        $pdo->exec("CREATE DATABASE $db_name_quoted");
    }
} catch (PDOException $e) {
    sendResponse(false, 'Database Connection Failed: ' . $e->getMessage());
}

// 2. Run Migrations on the new database
try {
    $dsn = "pgsql:host=$db_host;port=$db_port;dbname=$db_name";
    $pdo = new PDO($dsn, $db_user, $db_pass, [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]);

    $schema_path = __DIR__ . '/../backend/schema.sql';
    if (!file_exists($schema_path)) {
        sendResponse(false, 'Migration failed: schema.sql not found at ' . $schema_path);
    }
    
    $sql = file_get_contents($schema_path);
    $pdo->exec($sql);

    // Initial Org
    $org_id = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(16384, 20479), mt_rand(32768, 49151), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535));
    $pdo->prepare("INSERT INTO organizations (id, name, plan_tier) VALUES (?, ?, 'free')")->execute([$org_id, $org_name]);
    
    // Admin User - Generate secure random password
    $admin_password = bin2hex(random_bytes(16)); // 32 character hex string
    $hashed_password = password_hash($admin_password, PASSWORD_BCRYPT);
    if ($hashed_password === false) {
        sendResponse(false, 'Failed to hash admin password.');
    }
    
    $user_id = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(16384, 20479), mt_rand(32768, 49151), mt_rand(0, 65535), mt_rand(0, 65535), mt_rand(0, 65535));
    $pdo->prepare("INSERT INTO users (id, email, password_hash, full_name) VALUES (?, ?, ?, ?)")->execute([$user_id, $admin_email, $hashed_password, 'Administrator']);

} catch (PDOException $e) {
    sendResponse(false, 'Migration failed: ' . $e->getMessage());
}

// 3. Write .env File
$env_content = "DATABASE_URL=postgres://$db_user:$db_pass@$db_host:$db_port/$db_name?sslmode=disable\n";
$env_content .= "PORT=8080\n";
$env_content .= "NODE_ENV=production\n";
// Generate JWT secret key for .env
$jwt_secret = bin2hex(random_bytes(32));
$env_content .= "JWT_SECRET_KEY=$jwt_secret\n";

if (@file_put_contents(__DIR__ . '/../.env', $env_content) === false) {
    sendResponse(false, 'Installation succeeded but failed to write .env file. Please check folder permissions.');
}

// Return success with admin password
sendResponse(true, 'Cortex has been successfully installed and configured.', [
    'admin_password' => $admin_password,
    'admin_email' => $admin_email,
    'message' => 'IMPORTANT: Save your admin password. You will need it to log in.'
]);

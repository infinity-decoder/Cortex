# Cortex Environment Setup Script (PowerShell)
# This script adds the project-isolated Go compiler and the Laragon Node.js to your PATH.

$ProjectRoot = Get-Item $PSScriptRoot
$GoBinPath = Join-Path $ProjectRoot.FullName "tools\go\bin"
$NodePath = "C:\laragon\bin\nodejs\node-v22"

# 1. Update Current Session Path
if (Test-Path $GoBinPath) {
    $env:Path = "$GoBinPath;" + $env:Path
    Write-Host "✅ Added Go to session PATH: $GoBinPath" -ForegroundColor Green
} else {
    Write-Warning "Go binary path not found at $GoBinPath. Please ensure Phase 0 Go installation was successful."
}

if (Test-Path $NodePath) {
    $env:Path = "$NodePath;" + $env:Path
    Write-Host "✅ Added Node.js to session PATH: $NodePath" -ForegroundColor Green
} else {
    Write-Warning "Node.js path not found at $NodePath."
}

Write-Host "`nCortex environment is now active for this session." -ForegroundColor Cyan
Write-Host "You can now run 'go version' and 'npm --version' in this terminal.`n"

# Optional: Persistence Instructions
Write-Host "To make these changes persistent for your USER account, run the following command:" -ForegroundColor Yellow
Write-Host "[System.Environment]::SetEnvironmentVariable('Path', `"`$GoBinPath;`$NodePath;`" + [System.Environment]::GetEnvironmentVariable('Path', 'User'), 'User')" -ForegroundColor Gray

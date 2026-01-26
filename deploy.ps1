Write-Host "Building Cortex Frontend..."
cd frontend
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed"
    exit 1
}
cd ..

Write-Host "Deploying to Root..."
# Move index.html to app.html to avoid overwriting index.php
Move-Item -Path "frontend\out\index.html" -Destination "frontend\out\app.html" -Force

# Copy all contents to root
Copy-Item -Path "frontend\out\*" -Destination "." -Recurse -Force

Write-Host "Deployment Complete. Access at http://localhost/Cortex"

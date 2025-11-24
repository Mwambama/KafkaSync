# cleanup.ps1
# WARNING: This deletes all database data, logs, and downloaded files!

Write-Host "🛑 Stopping Docker containers..."
docker-compose down

Write-Host "🧹 Cleaning up consumer log..."
if (Test-Path "consumer.log") { 
    Remove-Item "consumer.log" 
    Write-Host "   - Deleted consumer.log"
}

Write-Host "🗑️  Deleting Database Data (Postgres)..."
if (Test-Path "postgres-data") { 
    Remove-Item -Recurse -Force "postgres-data" 
    Write-Host "   - Deleted postgres-data"
}

Write-Host "🗑️  Deleting Cloud Data (MinIO)..."
if (Test-Path "minio-data") { 
    Remove-Item -Recurse -Force "minio-data" 
    Write-Host "   - Deleted minio-data"
}

Write-Host "✨ Clearing local download folders..."
if (Test-Path "incompletes") { 
    Get-ChildItem -Path "incompletes" -Recurse | Remove-Item -Recurse -Force
}
if (Test-Path "completes") { 
    Get-ChildItem -Path "completes" -Recurse | Remove-Item -Recurse -Force
}

Write-Host "🚀 Restarting Infrastructure..."
docker-compose up -d

# Wait a moment for Postgres to initialize before reporting success
Start-Sleep -Seconds 5

Write-Host "✅ System is fresh and ready for a demo!"
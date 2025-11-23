# usage: .\generate_data.ps1 "my-new-file.txt"
# If no filename is provided, it defaults to "auto-generated-file.txt"
# This script creates a dummy file in the ./sftp-data/ directory
param (
    [string]$FileName = "auto-generated-file.txt"
)

$TargetFolder = ".\sftp-data"

# Ensure the folder exists
if (-not (Test-Path $TargetFolder)) {
    New-Item -ItemType Directory -Force -Path $TargetFolder | Out-Null
}

$FilePath = Join-Path $TargetFolder $FileName

# Create a file with some random dummy content (Timestamp + Random Number)
$DummyContent = "This is a test file generated at $(Get-Date). Random ID: $(Get-Random)"
Set-Content -Path $FilePath -Value $DummyContent

Write-Host "✅ Generated file: $FilePath"
Write-Host "👉 Now run the producer and enter: $FileName"
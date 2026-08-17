#!/usr/bin/env pwsh
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  Bishal Puja Sewa - Docker Setup" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

function Require-Command($Name, $InstallHint) {
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        Write-Host "[ERROR] $Name is not installed or not on PATH." -ForegroundColor Red
        Write-Host "        $InstallHint" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host "[1/5] Checking Docker..." -ForegroundColor Yellow
Require-Command "docker" "Install Docker Desktop, then reopen this terminal."
docker compose version | Out-Null
Write-Host "  [OK] Docker Compose is available" -ForegroundColor Green

Write-Host ""
Write-Host "[2/5] Preparing environment files..." -ForegroundColor Yellow
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "  [OK] Created .env" -ForegroundColor Green
} else {
    Write-Host "  [OK] .env already exists" -ForegroundColor Green
}
if (-not (Test-Path "backend\.env")) {
    Copy-Item "backend\.env.example" "backend\.env"
    Write-Host "  [OK] Created backend\.env for local non-Docker runs" -ForegroundColor Green
} else {
    Write-Host "  [OK] backend\.env already exists" -ForegroundColor Green
}

$envValues = @{}
Get-Content ".env" | ForEach-Object {
    $line = $_.Trim()
    if ($line -and -not $line.StartsWith("#") -and $line.Contains("=")) {
        $parts = $line.Split("=", 2)
        $envValues[$parts[0].Trim()] = $parts[1].Trim()
    }
}
$frontendPort = if ($envValues.ContainsKey("FRONTEND_HOST_PORT")) { $envValues["FRONTEND_HOST_PORT"] } else { "3000" }
$backendPort = if ($envValues.ContainsKey("BACKEND_HOST_PORT")) { $envValues["BACKEND_HOST_PORT"] } else { "8081" }

Write-Host ""
Write-Host "[3/5] Building and starting containers..." -ForegroundColor Yellow
docker compose up -d --build postgres redis backend frontend

Write-Host ""
Write-Host "[4/5] Waiting for backend health..." -ForegroundColor Yellow
$healthy = $false
for ($i = 1; $i -le 40; $i++) {
    try {
        $response = Invoke-WebRequest -UseBasicParsing -TimeoutSec 2 -Uri "http://127.0.0.1:$backendPort/health"
        if ($response.StatusCode -eq 200) {
            $healthy = $true
            break
        }
    } catch {
        Start-Sleep -Seconds 2
    }
}
if (-not $healthy) {
    Write-Host "[ERROR] Backend did not become healthy. Showing logs:" -ForegroundColor Red
    docker compose logs --tail=120 backend
    exit 1
}
Write-Host "  [OK] Backend is healthy" -ForegroundColor Green

Write-Host ""
Write-Host "[5/5] Seeding demo data..." -ForegroundColor Yellow
docker compose run --rm seed

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  Setup Complete" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "Frontend: http://localhost:$frontendPort" -ForegroundColor Yellow
Write-Host "Backend:  http://localhost:$backendPort" -ForegroundColor Yellow
Write-Host ""
Write-Host "Default accounts:" -ForegroundColor White
Write-Host "Admin:    admin@bishalpujasewa.com / AdminPass123!" -ForegroundColor Yellow
Write-Host "Pandit:   pandit.atri@example.com / PanditPass123!" -ForegroundColor Yellow
Write-Host "Customer: customer@example.com / CustomerPass123!" -ForegroundColor Yellow
Write-Host ""
Write-Host "Useful commands:" -ForegroundColor White
Write-Host "docker compose logs -f" -ForegroundColor Yellow
Write-Host "docker compose down" -ForegroundColor Yellow
Write-Host "docker compose down -v   # reset database too" -ForegroundColor Yellow
Write-Host ""

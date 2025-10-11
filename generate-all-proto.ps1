# PowerShell skripta za generisanje protobuf fajlova - OBA SERVISA

Write-Host "========================================" -ForegroundColor Magenta
Write-Host "  Generisanje svih protobuf fajlova" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta
Write-Host ""

# 1. Tour Service
Write-Host "[1/2] Tour Service..." -ForegroundColor Cyan
Set-Location -Path "$PSScriptRoot\services\tour-service"

if (!(Test-Path -Path "gen\pb-go")) {
    New-Item -ItemType Directory -Path "gen\pb-go" -Force
}

protoc --go_out=.\gen\pb-go --go_opt=paths=source_relative --go-grpc_out=.\gen\pb-go --go-grpc_opt=paths=source_relative .\proto\tour.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Tour Service - Protobuf generisan" -ForegroundColor Green
    go mod tidy | Out-Null
    Write-Host "  ✓ Tour Service - Dependencies ažurirani" -ForegroundColor Green
} else {
    Write-Host "  ✗ Tour Service - GREŠKA!" -ForegroundColor Red
}

Write-Host ""

# 2. API Gateway
Write-Host "[2/2] API Gateway..." -ForegroundColor Cyan
Set-Location -Path "$PSScriptRoot\services\api-gateway"

if (!(Test-Path -Path "gen\pb-go")) {
    New-Item -ItemType Directory -Path "gen\pb-go" -Force
}

protoc --go_out=.\gen\pb-go --go_opt=paths=source_relative --go-grpc_out=.\gen\pb-go --go-grpc_opt=paths=source_relative .\proto\tour.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ API Gateway - Protobuf generisan" -ForegroundColor Green
    go mod tidy | Out-Null
    Write-Host "  ✓ API Gateway - Dependencies ažurirani" -ForegroundColor Green
} else {
    Write-Host "  ✗ API Gateway - GREŠKA!" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Magenta
Write-Host "  Gotovo!" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta
Write-Host ""
Write-Host "Sledeći korak: Pokrenite servise" -ForegroundColor Yellow
Write-Host "  1. cd services\tour-service && go run .\cmd\api\main.go" -ForegroundColor White
Write-Host "  2. cd services\api-gateway && go run main.go" -ForegroundColor White
Write-Host ""

# Vrati se u root
Set-Location -Path "$PSScriptRoot"

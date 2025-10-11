Write-Host "========================================"
Write-Host "  Generisanje protobuf (Docker)"
Write-Host "========================================"
Write-Host ""

$dockerRunning = docker ps 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker nije pokrenut!" -ForegroundColor Red
    exit 1
}

Write-Host "Docker je aktivan" -ForegroundColor Green
Write-Host ""

Write-Host "[1/2] Tour Service..." -ForegroundColor Cyan
docker run --rm -v "${PWD}/services/tour-service:/workspace" -w /workspace namely/protoc:1.51_1 -I=proto --go_out=gen/pb-go/tour --go_opt=paths=source_relative --go-grpc_out=gen/pb-go/tour --go-grpc_opt=paths=source_relative tour.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "  Tour Service OK" -ForegroundColor Green
}

Write-Host ""
Write-Host "[2/2] API Gateway..." -ForegroundColor Cyan
docker run --rm -v "${PWD}/services/api-gateway:/workspace" -w /workspace namely/protoc:1.51_1 -I=proto --go_out=gen/pb-go/tour --go_opt=paths=source_relative --go-grpc_out=gen/pb-go/tour --go-grpc_opt=paths=source_relative tour.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "  API Gateway OK" -ForegroundColor Green
}

Write-Host ""
Write-Host "Go mod tidy..." -ForegroundColor Cyan
Set-Location services/tour-service
go mod tidy | Out-Null
Set-Location ../api-gateway
go mod tidy | Out-Null
Set-Location ../..

Write-Host ""
Write-Host "Gotovo!" -ForegroundColor Green
Write-Host "Protobuf fajlovi su uspesno generisani!" -ForegroundColor Green

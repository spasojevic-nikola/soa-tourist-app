# PowerShell skripta za generisanje protobuf fajlova - API Gateway

Write-Host "Generisanje protobuf fajlova za API Gateway..." -ForegroundColor Green

# Pozicioniranje u api-gateway direktorijum
Set-Location -Path "$PSScriptRoot\services\api-gateway"

# Kreiranje gen/pb-go foldera ako ne postoji
if (!(Test-Path -Path "gen\pb-go")) {
    New-Item -ItemType Directory -Path "gen\pb-go" -Force
    Write-Host "Kreiran gen/pb-go folder" -ForegroundColor Yellow
}

# Generisanje protobuf fajlova
Write-Host "Pokretanje protoc komande..." -ForegroundColor Cyan
protoc --go_out=.\gen\pb-go --go_opt=paths=source_relative --go-grpc_out=.\gen\pb-go --go-grpc_opt=paths=source_relative .\proto\tour.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Protobuf fajlovi uspešno generisani!" -ForegroundColor Green
    
    # Listanje generisanih fajlova
    Write-Host "`nGenerisani fajlovi:" -ForegroundColor Cyan
    Get-ChildItem -Path "gen\pb-go\tour*.pb.go" | ForEach-Object { Write-Host "  - $($_.Name)" -ForegroundColor White }
    
    # go mod tidy
    Write-Host "`nPokretanje go mod tidy..." -ForegroundColor Cyan
    go mod tidy
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Dependencies ažurirani!" -ForegroundColor Green
    } else {
        Write-Host "✗ Greška pri ažuriranju dependencies!" -ForegroundColor Red
    }
} else {
    Write-Host "✗ Greška pri generisanju protobuf fajlova!" -ForegroundColor Red
    Write-Host "Proverite da li je protoc instaliran: protoc --version" -ForegroundColor Yellow
}

Write-Host "`nGotovo!" -ForegroundColor Green

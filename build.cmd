@echo off
setlocal
cd /d "%~dp0"
if not exist "bin" mkdir "bin"
set GOOS=windows
set GOARCH=amd64
go build -trimpath -ldflags="-s -w" -o "bin\touch.exe" ./cmd/touch
if errorlevel 1 exit /b 1
bin\touch.exe --version
exit /b 0

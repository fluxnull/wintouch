@echo off
setlocal
set "DRV=%~d0"
set "OUT=%DRV%\go-workspace\bin\touch"
if not exist "%OUT%" mkdir "%OUT%"
copy /Y "%~dp0bin\touch.exe" "%OUT%\touch.exe"
exit /b %ERRORLEVEL%

@echo off
setlocal

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0setup.ps1"
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Setup failed. Check the error above.
    pause
    exit /b %ERRORLEVEL%
)

pause

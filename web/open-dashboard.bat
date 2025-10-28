@echo off
echo ================================
echo RabbitMQ Monitoring Dashboard
echo ================================
echo.
echo Starting dashboard...
echo.
echo Dashboard URL: http://localhost:3000/dashboard.html
echo API URL: http://localhost:8080/api/v1/monitoring
echo.
echo Press Ctrl+C to stop the server
echo.

cd /d "%~dp0"

REM Check if Python is available
python --version >nul 2>&1
if %errorlevel% equ 0 (
    echo Using Python HTTP Server...
    echo.
    start http://localhost:3000/dashboard.html
    python -m http.server 3000
) else (
    echo Python not found. Opening dashboard directly...
    echo Note: If you see CORS errors, install Python and run this script again.
    echo.
    start dashboard.html
)

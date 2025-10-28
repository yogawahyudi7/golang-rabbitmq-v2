#!/bin/bash

echo "================================"
echo "RabbitMQ Monitoring Dashboard"
echo "================================"
echo ""
echo "Starting dashboard..."
echo ""
echo "Dashboard URL: http://localhost:3000/dashboard.html"
echo "API URL: http://localhost:8080/api/v1/monitoring"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

cd "$(dirname "$0")"

# Check if Python is available
if command -v python3 &> /dev/null; then
    echo "Using Python3 HTTP Server..."
    echo ""

    # Open browser
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        open http://localhost:3000/dashboard.html
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        xdg-open http://localhost:3000/dashboard.html 2>/dev/null || \
        firefox http://localhost:3000/dashboard.html 2>/dev/null || \
        google-chrome http://localhost:3000/dashboard.html 2>/dev/null
    fi

    python3 -m http.server 3000
elif command -v python &> /dev/null; then
    echo "Using Python HTTP Server..."
    echo ""

    # Open browser
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open http://localhost:3000/dashboard.html
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        xdg-open http://localhost:3000/dashboard.html 2>/dev/null || \
        firefox http://localhost:3000/dashboard.html 2>/dev/null || \
        google-chrome http://localhost:3000/dashboard.html 2>/dev/null
    fi

    python -m SimpleHTTPServer 3000
else
    echo "Python not found. Opening dashboard directly..."
    echo "Note: If you see CORS errors, install Python and run this script again."
    echo ""

    if [[ "$OSTYPE" == "darwin"* ]]; then
        open dashboard.html
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        xdg-open dashboard.html 2>/dev/null || \
        firefox dashboard.html 2>/dev/null || \
        google-chrome dashboard.html 2>/dev/null
    fi
fi

@echo off
echo Starting WebhookRelay...
echo.

REM Check if .env exists
if not exist ".env" (
    echo ERROR: .env file not found!
    echo Please copy .env.example to .env and configure your Supabase credentials.
    echo.
    pause
    exit /b 1
)

echo Starting backend API server...
start "WebhookRelay API" cmd /k "go run cmd/api/main.go"

timeout /t 3

echo Starting mock receiver...
start "Mock Receiver" cmd /k "go run cmd/mock/main.go"

timeout /t 2

echo Starting frontend...
cd frontend
start "Frontend Dev Server" cmd /k "npm run dev"
cd ..

echo.
echo All services started!
echo - API: http://localhost:8080
echo - Mock: http://localhost:9090
echo - Frontend: http://localhost:5173
echo.
pause

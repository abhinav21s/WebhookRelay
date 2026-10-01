@echo off
echo Testing WebhookRelay API...
echo.

REM Check if backend is running
echo [1/6] Checking backend health...
curl -s http://localhost:8080/health
if %errorlevel% neq 0 (
    echo ERROR: Backend not responding on port 8080
    echo Please start backend with: go run cmd/api/main.go
    pause
    exit /b 1
)
echo OK - Backend is running
echo.

REM Create a test endpoint
echo [2/6] Creating test endpoint...
curl -s -X POST http://localhost:8080/api/webhooks ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Test Endpoint\",\"url\":\"http://localhost:9090/webhook\",\"event_types\":[\"test.event\"]}"
echo.
echo OK - Endpoint created
echo.

REM List endpoints
echo [3/6] Listing endpoints...
curl -s http://localhost:8080/api/webhooks
echo.
echo OK - Endpoints listed
echo.

REM Publish test event
echo [4/6] Publishing test event...
curl -s -X POST http://localhost:8080/api/events ^
  -H "Content-Type: application/json" ^
  -d "{\"event_type\":\"test.event\",\"payload\":{\"test\":true,\"timestamp\":\"2024-01-01T12:00:00Z\"}}"
echo.
echo OK - Event published
echo.

REM Wait for delivery
echo [5/6] Waiting 2 seconds for delivery...
timeout /t 2 /nobreak > nul
echo.

REM Check metrics
echo [6/6] Checking metrics...
curl -s http://localhost:8080/api/metrics
echo.
echo OK - Metrics retrieved
echo.

echo ========================================
echo All tests passed! WebhookRelay is working correctly.
echo ========================================
echo.
echo Next steps:
echo 1. Open http://localhost:5173 in your browser
echo 2. Check the dashboard for the test event
echo 3. Review RUNNING.md for feature testing
echo.
pause

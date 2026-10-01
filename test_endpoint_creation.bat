@echo off
echo Testing Endpoint Creation...
echo.

echo Step 1: Check backend health
curl -s http://localhost:8080/health
echo.
echo.

echo Step 2: Create test endpoint
curl -s -X POST http://localhost:8080/api/webhooks ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"CLI Test Endpoint\",\"url\":\"http://localhost:9090/webhook\",\"event_types\":[\"test.event\"]}"
echo.
echo.

echo Step 3: List all endpoints
curl -s http://localhost:8080/api/webhooks
echo.
echo.

echo ========================================
echo If you see endpoint JSON above, backend is working!
echo Now check the UI in browser: http://localhost:5173
echo ========================================
echo.
pause

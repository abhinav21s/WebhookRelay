# Fix: Retry Logic Not Working

## Problem
When you set failure mode to "500 Internal Server Error" and send an event, it shows "Delivered" instead of retrying with failures.

## Root Causes

### 1. Mock Receiver Confusion
You might be running **two separate mock receivers**:
- Standalone: `go run cmd/mock/main.go` (port 9090)
- Built-in API: Part of main API server (port 8080)

They don't share state!

### 2. Endpoint URL Mismatch
Your endpoints might point to the wrong mock receiver.

## ✅ Solution

### Option A: Use Built-in Mock (RECOMMENDED)

**1. Stop the separate mock receiver**
- If you're running `cmd/mock/main.go` in Terminal 2, stop it (Ctrl+C)

**2. Create endpoints pointing to API server**
```
URL: http://localhost:8080/webhook
```

**3. Control mode via API server**
The failure mode selector in UI will call:
```
POST http://localhost:8080/api/mock/control
```

### Option B: Use Standalone Mock

**1. Keep mock receiver running**
```bash
go run cmd/mock/main.go
```

**2. Create endpoints pointing to mock**
```
URL: http://localhost:9090/webhook
```

**3. Control mode manually**
```bash
curl -X POST http://localhost:9090/control \
  -H "Content-Type: application/json" \
  -d '{"mode":"error500"}'
```

**Problem:** The UI failure mode selector won't work because it calls port 8080, not 9090!

---

## 🔧 Complete Fix Steps

### Step 1: Stop Separate Mock (if running)

If you have Terminal 2 running `cmd/mock/main.go`:
- Press Ctrl+C to stop it

### Step 2: Delete Old Endpoints

In browser:
1. Go to Endpoints page
2. Delete all existing endpoints that point to `localhost:9090`

### Step 3: Create New Endpoint

1. Click "New Endpoint"
2. Fill in:
   - Name: `Test Retry`
   - URL: `http://localhost:8080/webhook`  ← **IMPORTANT: Port 8080, not 9090!**
   - Event Types: `test.retry`
3. Click "Create"

### Step 4: Test Failure Mode

1. Click on the new endpoint
2. Set Failure Mode to **"500 Internal Server Error"**
3. Click "Send Test Event"
4. Go to Events page
5. Click on the event
6. **Watch the delivery timeline!**

You should see:
```
● Attempt 1  Failed (500)   ~100ms
  (wait 1 second)
● Attempt 2  Failed (500)   ~100ms
  (wait 2 seconds)
● Attempt 3  Failed (500)   ~100ms
  (wait 4 seconds)
● Attempt 4  Failed (500)   ~100ms
  (wait 8 seconds)
● Attempt 5  Failed (500)   ~100ms
  
Status: Dead Letter (red)
```

### Step 5: Check Terminal 1 (Backend)

You should see logs like:
```
[Mock Receiver] Received webhook - Event: xxx
[Mock Receiver] Responding with 500 Internal Server Error
Worker 1: Attempt 1 failed (status 500), retrying in 1s (attempt 2/5)
Worker 1: Attempt 2 failed (status 500), retrying in 2s (attempt 3/5)
...
Worker 1: Max attempts reached for event xxx, moving to dead letter
```

---

## 🐛 If Still Not Working

### Debug Checklist:

**1. Check which mock is being used**

In Terminal 1, when you send event, look for:
```
[Mock Receiver] Received webhook
```

If you see this = built-in mock is working ✅

**2. Check the failure mode was set**

```bash
curl http://localhost:8080/api/mock/mode
```

Should return:
```json
{"mode":"error500"}
```

**3. Verify endpoint URL**

In browser, check endpoint detail:
```
URL: http://localhost:8080/webhook
```

NOT:
```
URL: http://localhost:9090/webhook  ❌
```

**4. Check event status updates**

Browser console should show:
```
🔍 Fetching events from: http://localhost:8080/api/events
```

Every 2 seconds.

**5. Check delivery attempts are logged**

In Supabase SQL Editor:
```sql
SELECT 
  event_id, 
  attempt_number, 
  status_code, 
  error_message,
  created_at
FROM delivery_attempts
ORDER BY created_at DESC
LIMIT 10;
```

Should show multiple attempts with status_code = 500

---

## 🎯 Quick Test Script

Run this to verify everything:

```bash
# 1. Set mock to error mode
curl -X POST http://localhost:8080/api/mock/control \
  -H "Content-Type: application/json" \
  -d '{"mode":"error500"}'

# 2. Verify mode was set
curl http://localhost:8080/api/mock/mode

# 3. Create test event
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{"event_type":"test.retry","payload":{"test":true}}'

# 4. Wait 20 seconds for all retries
timeout /t 20

# 5. Check delivery attempts
curl http://localhost:8080/api/deliveries
```

---

## Expected Terminal 1 Output

```
POST /api/mock/control - 200 (5ms)
[Mock Receiver] Mode changed to: error500

GET /api/mock/mode - 200 (2ms)

POST /api/events - 201 (45ms)
Worker 1: Processing attempt 1/5 for event xxx to http://localhost:8080/webhook
[Mock Receiver] Received webhook - Event: xxx
[Mock Receiver] Responding with 500 Internal Server Error
Worker 1: Attempt 1 failed (status 500), retrying in 1s (attempt 2/5)

(1 second later)
Worker 1: Processing attempt 2/5 for event xxx
[Mock Receiver] Received webhook - Event: xxx
[Mock Receiver] Responding with 500 Internal Server Error
Worker 1: Attempt 2 failed (status 500), retrying in 2s (attempt 3/5)

(2 seconds later)
Worker 1: Processing attempt 3/5 for event xxx
[Mock Receiver] Received webhook - Event: xxx
[Mock Receiver] Responding with 500 Internal Server Error
Worker 1: Attempt 3 failed (status 500), retrying in 4s (attempt 4/5)

(4 seconds later)
Worker 1: Processing attempt 4/5 for event xxx
[Mock Receiver] Received webhook - Event: xxx
[Mock Receiver] Responding with 500 Internal Server Error
Worker 1: Attempt 4 failed (status 500), retrying in 8s (attempt 5/5)

(8 seconds later)
Worker 1: Processing attempt 5/5 for event xxx
[Mock Receiver] Received webhook - Event: xxx
[Mock Receiver] Responding with 500 Internal Server Error
Worker 1: Max attempts reached for event xxx, moving to dead letter
```

---

## Summary

**The key issue:** Using wrong mock receiver URL.

**The fix:** Use `http://localhost:8080/webhook` (built-in mock) instead of `http://localhost:9090/webhook` (separate mock).

This way, the UI failure mode selector actually controls the mock your endpoint calls! 🚀

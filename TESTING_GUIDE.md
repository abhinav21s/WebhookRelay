# WebhookRelay - Step-by-Step Testing Guide

Complete walkthrough to test every feature in 30 minutes.

## Prerequisites Setup

### 1. Start All Services

Open **3 separate terminals**:

**Terminal 1 - API Server:**
```bash
go run cmd/api/main.go
```
✅ Wait for: "Worker pool started with 10 workers"

**Terminal 2 - Mock Receiver:**
```bash
go run cmd/mock/main.go
```
✅ Wait for: "Mock receiver started on :9090"

**Terminal 3 - Frontend:**
```bash
cd frontend
npm run dev
```
✅ Wait for: "Local: http://localhost:5173/"

### 2. Open Dashboard

Open browser: **http://localhost:5173**

You should see the dark-themed dashboard.

---

## Feature Testing Checklist

### ✅ TEST 1: Dashboard Overview (2 minutes)

**What to check:**
1. Navigate to **Dashboard** page
2. Look at the 4 metric cards at top:
   - Total Events (24h)
   - Success Rate
   - Avg Latency
   - Active Endpoints
3. Scroll down to see:
   - Success Rate Chart
   - Endpoint Health (empty for now)
   - Recent Events (empty for now)

**Expected Result:**
- All metrics show 0 (no data yet)
- UI loads without errors
- Dark theme looks good

---

### ✅ TEST 2: Create Endpoint (3 minutes)

**Steps:**
1. Click **"Endpoints"** in top navigation
2. Click **"New Endpoint"** button (top right)
3. Fill in the form:
   - **Name:** `Test Endpoint`
   - **URL:** `http://localhost:9090/webhook`
   - **Event Types:** `test.event, payment.completed`
   - **Secret:** (leave empty - auto-generates)
4. Click **"Create"**

**Expected Result:**
- ✅ Modal closes
- ✅ New endpoint appears in table
- ✅ Status badge shows "Active" (green)
- ✅ Event types displayed as pills
- ✅ URL shown in monospace font

**Dashboard Check:**
- Go back to Dashboard
- "Active Endpoints" should now show **1**

---

### ✅ TEST 3: Send Successful Test Event (5 minutes)

**Steps:**
1. Go to **Endpoints** page
2. Click on **"Test Endpoint"** (the name)
3. You're now on Endpoint Detail page
4. Find **"Failure Mode"** dropdown
5. Make sure it's set to **"Success (200)"**
6. Click **"Send Test Event"** button
7. A modal opens with:
   - Event Type: `test.event` (selected)
   - Payload: JSON (pre-filled)
8. Click **"Send Event"**

**Expected Result:**
- ✅ Modal closes
- ✅ Toast/success message (optional)

**Verify Delivery:**
1. Click **"Events"** in navigation
2. You should see your event in the list
3. Status badge: "Pending" → "Delivering" → "Delivered" (green)
4. Click on the event ID to open Event Detail page

**Event Detail Page:**
- ✅ Event Type: `test.event`
- ✅ Status: "Delivered" (green badge)
- ✅ Payload shown in formatted JSON
- ✅ Delivery Timeline shows **1 attempt**:
  - Attempt 1
  - Timestamp
  - Status: "Succeeded (200)"
  - Latency: ~100-200ms

**Mock Receiver Check:**
- Go to Terminal 2 (mock receiver)
- You should see log: `[Mock Receiver] Received webhook`
- Headers logged: Event ID, Signature, Timestamp

---

### ✅ TEST 4: Retry Logic - The Star Feature! (10 minutes)

This is the most impressive feature to demonstrate!

**Steps:**
1. Go back to **Endpoint Detail** page
2. Change **"Failure Mode"** to **"500 Internal Server Error"**
3. Click **"Send Test Event"**
4. Click **"Send Event"**
5. Immediately go to **Events** page
6. Click on the new event to open Event Detail

**Watch the Magic:**
- Status changes from "Pending" → "Delivering" → "Retrying"
- Delivery Timeline updates in real-time (every 2 seconds)
- Watch attempts appear one by one:

```
● Attempt 1  12:00:00  Failed (500)   150ms
  (immediately)

● Attempt 2  12:00:01  Failed (500)   140ms
  (+1 second delay)

● Attempt 3  12:00:03  Failed (500)   135ms
  (+2 seconds delay)

● Attempt 4  12:00:07  Failed (500)   142ms
  (+4 seconds delay)

● Attempt 5  12:00:15  Failed (500)   138ms
  (+8 seconds delay)
```

**Expected Result:**
- ✅ 5 attempts total
- ✅ Exponential backoff visible: 0s, +1s, +2s, +4s, +8s
- ✅ Each attempt shows "Failed (500)"
- ✅ After attempt 5, status changes to **"Dead Letter"** (purple/red)
- ✅ Timeline animates as new attempts appear
- ✅ No more retries after 5th attempt

**Mock Receiver Check:**
- Terminal 2 shows 5 webhook requests
- Each logged with "Responding with 500"

---

### ✅ TEST 5: Dead Letter Queue (2 minutes)

**Steps:**
1. Click **"Dead Letters"** in navigation
2. You should see the failed event from Test 4

**Expected Result:**
- ✅ Event listed with failed timestamp
- ✅ "Replay" button visible
- ✅ "Inspect" button navigates to Event Detail

---

### ✅ TEST 6: Replay Event (3 minutes)

**Steps:**
1. From Dead Letters page
2. **IMPORTANT:** First go back to Endpoint Detail
3. Change **"Failure Mode"** back to **"Success (200)"**
4. Return to **Dead Letters** page
5. Click **"Replay"** button on the failed event
6. Confirm if prompted

**Expected Result:**
- ✅ Navigates to new event detail page (different event ID)
- ✅ Same payload as original
- ✅ Status: "Delivered" (green)
- ✅ Timeline shows 1 successful attempt
- ✅ Original event still in Dead Letters (history preserved)

---

### ✅ TEST 7: Multiple Failure Modes (5 minutes)

Test each failure mode to see different behaviors:

#### A. Timeout Mode
1. Endpoint Detail → Failure Mode: **"Timeout"**
2. Send test event
3. Watch Event Detail
4. **Expected:** Error message shows "timeout" or "context deadline exceeded"
5. Retries happen (timeout is retryable)

#### B. Rate Limit (429)
1. Endpoint Detail → Failure Mode: **"429 Rate Limited"**
2. Send test event
3. **Expected:** Status code 429, retries happen

#### C. Bad Request (400)
1. Endpoint Detail → Failure Mode: **"400 Bad Request"**
2. Send test event
3. **Expected:** Status code 400, **NO RETRIES** (4xx errors are not retryable)
4. Immediately marked as "Failed"

---

### ✅ TEST 8: Multiple Endpoints (Fan-out) (5 minutes)

**Steps:**
1. Create 2 more endpoints:
   
   **Endpoint 2:**
   - Name: `Endpoint A`
   - URL: `http://localhost:9090/webhook-a`
   - Event Types: `multi.test`
   
   **Endpoint 3:**
   - Name: `Endpoint B`
   - URL: `http://localhost:9090/webhook-b`
   - Event Types: `multi.test`

2. Set **all endpoints** to Success mode

3. Go to any endpoint detail with `multi.test` event type

4. Send test event with type: `multi.test`

5. Go to Events page and click the event

**Expected Result:**
- ✅ Single event shows **multiple delivery timelines** (one per endpoint)
- ✅ All endpoints receive the webhook
- ✅ Each has independent delivery status
- ✅ Mock receiver logs show 3 separate requests

---

### ✅ TEST 9: Enable/Disable Endpoint (2 minutes)

**Steps:**
1. Go to **Endpoints** page
2. Click the **Active** badge on an endpoint
3. It should toggle to **"Disabled"** (gray)
4. Send a test event that matches this endpoint's event types
5. Check Events page

**Expected Result:**
- ✅ Event is created
- ✅ But NOT delivered to disabled endpoint
- ✅ Only active endpoints receive it

---

### ✅ TEST 10: Metrics Update (2 minutes)

**Steps:**
1. After running several tests above
2. Go to **Dashboard**
3. Refresh page if needed

**Expected Result:**
- ✅ Total Events (24h): Shows count of all events
- ✅ Success Rate: Shows percentage (e.g., 60% if some failed)
- ✅ Avg Latency: Shows average delivery time in ms
- ✅ Active Endpoints: Shows count of enabled endpoints
- ✅ Success Rate Chart: Shows trend line
- ✅ Recent Events: Lists last 20 events
- ✅ Endpoint Health: Green dots for active endpoints

---

### ✅ TEST 11: Real-time Updates (2 minutes)

**Steps:**
1. Set an endpoint to "500" failure mode
2. Send test event
3. Open Event Detail page
4. **DON'T REFRESH** - just watch

**Expected Result:**
- ✅ Status badge updates automatically (every 2s)
- ✅ New attempts appear in timeline without refresh
- ✅ Attempt numbers increment
- ✅ Status changes from "Retrying (1/5)" → "Retrying (2/5)" → etc.

---

### ✅ TEST 12: API Testing with cURL (Optional)

If you want to test the API directly:

**Create Endpoint:**
```bash
curl -X POST http://localhost:8080/api/webhooks \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"API Test\",\"url\":\"http://localhost:9090/webhook\",\"event_types\":[\"api.test\"]}"
```

**Publish Event:**
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d "{\"event_type\":\"api.test\",\"payload\":{\"test\":true}}"
```

**Check Metrics:**
```bash
curl http://localhost:8080/api/metrics
```

**Control Mock Mode:**
```bash
curl -X POST http://localhost:9090/control \
  -H "Content-Type: application/json" \
  -d "{\"mode\":\"error500\"}"
```

---

## Quick Demo Script (5 Minutes)

For showing someone the app:

### 1. Show Dashboard (30s)
"Here's our webhook delivery system dashboard with metrics"

### 2. Create Endpoint (30s)
"Let me create an endpoint - this is where webhooks get sent"

### 3. Successful Delivery (1min)
"Send a test event... see? Delivered in ~100ms, status is green"

### 4. Demonstrate Retry Logic (2min) ⭐
"Now watch this - I'll make it fail...
- First attempt fails immediately
- Retry after 1 second
- Retry after 2 more seconds
- Retry after 4 more seconds
- Retry after 8 more seconds
- After 5 attempts, it goes to dead letter queue"

"This is exponential backoff - it prevents overwhelming a struggling server"

### 5. Replay (1min)
"Now I can switch it back to success mode and replay the failed event - see? Works now!"

---

## Troubleshooting

### Event Not Delivering
- ✅ Check endpoint is Active (green)
- ✅ Check event type matches endpoint subscription
- ✅ Check mock receiver is running
- ✅ Check backend logs in Terminal 1

### Timeline Not Updating
- ✅ Check browser console for errors (F12)
- ✅ Verify backend is responding: http://localhost:8080/health
- ✅ Refresh the page

### Mock Not Responding
- ✅ Check Terminal 2 is running
- ✅ Try: `curl http://localhost:9090/webhook`
- ✅ Restart mock receiver

---

## What You've Tested

After completing this guide, you've verified:

✅ Endpoint CRUD operations
✅ Event publishing
✅ Successful webhook delivery
✅ Exponential backoff retry logic
✅ All 5 failure modes
✅ Dead letter queue
✅ Event replay
✅ Multi-endpoint fan-out
✅ Enable/disable endpoints
✅ Real-time UI updates
✅ Metrics aggregation
✅ HMAC signature generation
✅ Timeline visualization

**You're now ready to demo or deploy! 🚀**

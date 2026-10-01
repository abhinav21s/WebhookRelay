# Running & Testing WebhookRelay

Complete guide to run the application and test all features.

## Starting the Application

### Start Backend

```bash
# From project root
go run cmd/api/main.go
```

Expected output:
```
2024/01/01 12:00:00 Starting WebhookRelay API on :8080
2024/01/01 12:00:00 Connected to Supabase successfully
2024/01/01 12:00:00 Worker pool started with 10 workers
```

### Start Mock Receiver (for testing)

In a separate terminal:

```bash
go run cmd/mock/main.go
```

Expected output:
```
2024/01/01 12:00:00 Mock receiver started on :9090
2024/01/01 12:00:00 Default mode: success
```

### Start Frontend

In a separate terminal:

```bash
cd frontend
npm run dev
```

Expected output:
```
  VITE v5.x.x  ready in xxx ms
  ➜  Local:   http://localhost:5173/
```

Open your browser to: **http://localhost:5173**

## Testing Core Features

### Feature 1: Register an Endpoint

**Via API:**
```bash
curl -X POST http://localhost:8080/api/webhooks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Payment Service",
    "url": "http://localhost:9090/webhook",
    "event_types": ["payment.completed", "payment.failed"],
    "secret": "my-secret-key"
  }'
```

**Via Dashboard:**
1. Navigate to **Endpoints** page
2. Click **New Endpoint** button
3. Fill in the form:
   - Name: Payment Service
   - URL: http://localhost:9090/webhook
   - Event Types: payment.completed, payment.failed
   - Secret: my-secret-key
4. Click **Create**

**What to Check:**
- ✅ Endpoint appears in the list
- ✅ Status badge shows "Active" (green)
- ✅ Event types are displayed correctly

### Feature 2: Publish an Event

**Via API:**
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "payment.completed",
    "payload": {
      "order_id": "ORD-12345",
      "amount": 99.99,
      "currency": "USD"
    }
  }'
```

**Via Dashboard:**
1. Go to **Endpoint Detail** page
2. Click **Send Test Event** button
3. Select event type
4. (Optional) Edit payload
5. Click **Send**

**What to Check:**
- ✅ Event appears in Events list immediately
- ✅ Status changes from "Pending" → "Delivering" → "Delivered"
- ✅ Success badge (green) appears
- ✅ Delivery attempt logged with ~100-200ms latency

### Feature 3: Retry Logic (The Star Feature!)

This demonstrates the exponential backoff retry mechanism.

**Setup:**
1. Go to **Endpoint Detail** page
2. Find **Failure Mode Selector** (dropdown)
3. Select **"500 Internal Server Error"**
4. Click **Send Test Event**

**What to Check:**
- ✅ Event status shows "Retrying"
- ✅ Delivery timeline shows multiple attempts:
  - Attempt 1: Failed (500) - ~0s
  - Attempt 2: Failed (500) - ~1s later
  - Attempt 3: Failed (500) - ~2s later
  - Attempt 4: Failed (500) - ~4s later
  - Attempt 5: Failed (500) - ~8s later
- ✅ After 5 attempts: Status changes to "Dead Letter" (red)
- ✅ Each attempt visible in real-time (polling updates)

**Other Failure Modes to Test:**

**Timeout:**
```
1. Set mode to "Timeout"
2. Send test event
3. Watch retries with "timeout" error message
```

**Rate Limit (429):**
```
1. Set mode to "429 Rate Limited"
2. Send test event
3. Watch retries with 429 status code
```

**Client Error (400):**
```
1. Set mode to "400 Bad Request"
2. Send test event
3. Note: No retries (4xx errors are not retryable)
4. Immediately marked as failed
```

### Feature 4: HMAC Signature Verification

**Check Outgoing Webhook Headers:**

Inspect mock receiver logs to see headers:
```
X-Webhook-Signature: sha256=abc123...
X-Webhook-ID: event-uuid
X-Webhook-Timestamp: 1234567890
Content-Type: application/json
```

**Verify Signature (example in Node.js):**
```javascript
const crypto = require('crypto');

function verifySignature(payload, signature, secret) {
  const hmac = crypto.createHmac('sha256', secret);
  hmac.update(payload);
  const computed = 'sha256=' + hmac.digest('hex');
  return computed === signature;
}
```

### Feature 5: Live Status Updates

**Test Real-Time Polling:**
1. Send a test event with "500" failure mode
2. Keep **Event Detail** page open
3. Watch status badges update automatically every 2-3 seconds
4. See new attempts appear in timeline without page refresh

**What to Check:**
- ✅ Status badge updates without refresh
- ✅ Timeline items appear automatically
- ✅ Latency values update in real-time
- ✅ Overall status progresses: Pending → Delivering → Retrying → Dead Letter

### Feature 6: Event Replay

**Steps:**
1. Go to **Dead Letters** page
2. Find a permanently failed event
3. Click **Replay** button
4. Confirmation dialog appears
5. Click **Confirm**

**Before Replay:**
- Set failure mode back to **"Success"** in mock receiver

**What to Check:**
- ✅ New event created with same payload
- ✅ Event re-queued and delivered successfully
- ✅ New delivery timeline created
- ✅ Original event remains in Dead Letters (history preserved)

### Feature 7: Endpoint Management

**Disable an Endpoint:**
1. Go to **Endpoints** page
2. Toggle **Active** switch to OFF
3. Send a test event matching that endpoint's event types

**What to Check:**
- ✅ Status badge changes to "Disabled" (gray)
- ✅ Event is published but NOT delivered to disabled endpoint
- ✅ Other active endpoints still receive the event

**Edit Endpoint:**
1. Click endpoint name to open detail page
2. Click **Edit** button
3. Change event types or URL
4. Click **Save**

**Delete Endpoint:**
1. Click **Delete** button
2. Confirm deletion
3. Endpoint removed from list (soft delete - history preserved)

### Feature 8: Metrics & Dashboard

**Navigate to Overview/Dashboard:**

**What to Check:**
- ✅ **Total Events (24h)**: Count of recent events
- ✅ **Success Rate**: Percentage of successful deliveries
- ✅ **Avg Latency**: Average delivery time in ms
- ✅ **Active Endpoints**: Count of enabled endpoints
- ✅ **Success Rate Chart**: Line graph showing trends over 24 hours
- ✅ **Recent Events**: Live feed of latest 20 events
- ✅ **Endpoint Health**: Visual indicators (green/yellow/red dots)

### Feature 9: Delivery Timeline (Most Important UI)

**Steps:**
1. Create event with multiple retries (use 500 failure mode)
2. Navigate to **Event Detail** page
3. Observe the vertical timeline

**What to Check:**
- ✅ Timeline shows all attempts in chronological order
- ✅ Each attempt displays:
  - Attempt number
  - Timestamp
  - Status code or error message
  - Latency
  - Color-coded status (red for failed, green for success)
- ✅ Timeline updates live as new attempts occur
- ✅ Clear visual progression of retry delays

### Feature 10: Multiple Endpoints Fan-out

**Setup:**
1. Create 3 different endpoints:
   - Endpoint A: http://localhost:9090/webhook-a (success mode)
   - Endpoint B: http://localhost:9090/webhook-b (500 mode)
   - Endpoint C: http://localhost:9090/webhook-c (success mode)
2. All subscribed to "multi.test" event type

**Test:**
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "multi.test",
    "payload": {"test": "fan-out"}
  }'
```

**What to Check:**
- ✅ Single event creates 3 separate delivery jobs
- ✅ Endpoint A: Delivers successfully
- ✅ Endpoint B: Retries with exponential backoff
- ✅ Endpoint C: Delivers successfully
- ✅ Each endpoint shows independent delivery timeline
- ✅ Worker pool processes all concurrently

## API Testing with cURL

### List All Endpoints
```bash
curl http://localhost:8080/api/webhooks
```

### Get Endpoint by ID
```bash
curl http://localhost:8080/api/webhooks/{id}
```

### Update Endpoint
```bash
curl -X PATCH http://localhost:8080/api/webhooks/{id} \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'
```

### Delete Endpoint
```bash
curl -X DELETE http://localhost:8080/api/webhooks/{id}
```

### List Events
```bash
curl http://localhost:8080/api/events
```

### Get Event Detail
```bash
curl http://localhost:8080/api/events/{id}
```

### List Deliveries for Event
```bash
curl "http://localhost:8080/api/deliveries?event_id={id}"
```

### Get Metrics
```bash
curl http://localhost:8080/api/metrics
```

### Control Mock Receiver Mode
```bash
# Set to return 500
curl -X POST http://localhost:9090/control \
  -H "Content-Type: application/json" \
  -d '{"mode": "error500"}'

# Set to timeout
curl -X POST http://localhost:9090/control \
  -H "Content-Type: application/json" \
  -d '{"mode": "timeout"}'

# Set to success
curl -X POST http://localhost:9090/control \
  -H "Content-Type: application/json" \
  -d '{"mode": "success"}'
```

## Performance Testing

### Load Test with Multiple Events
```bash
# Send 100 events rapidly
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/events \
    -H "Content-Type: application/json" \
    -d "{\"event_type\":\"load.test\",\"payload\":{\"index\":$i}}" &
done
```

**What to Check:**
- ✅ All events queued without blocking
- ✅ Worker pool processes concurrently (max 10 at a time)
- ✅ No dropped events
- ✅ Metrics update correctly

## Troubleshooting

### Events Not Delivering
- Check if endpoint is active (enabled)
- Verify mock receiver is running
- Check event type matches endpoint subscription
- Inspect backend logs for errors

### Timeline Not Updating
- Verify React Query polling is working (check browser console)
- Ensure backend is returning delivery attempts
- Check network tab for API calls

### Signature Verification Failing
- Ensure secret matches between endpoint and verification code
- Check payload is used as-is (no modifications)
- Verify HMAC algorithm is SHA256

## Demo Script for Interviewers

**5-Minute Live Demo:**

1. **Show Dashboard** (30s)
   - Overview with metrics
   - Active endpoints list

2. **Create Endpoint** (30s)
   - Quick form fill
   - Subscribe to event type

3. **Send Successful Event** (1min)
   - Test event delivers immediately
   - Show timeline with ~100ms latency
   - Point out HMAC signature

4. **Demonstrate Retry Logic** (2min) ⭐
   - Set failure mode to "500"
   - Send test event
   - Watch live timeline update with retries
   - Point out exponential backoff delays (1s, 2s, 4s, 8s...)
   - Show final Dead Letter status

5. **Replay Event** (1min)
   - Switch mode back to success
   - Replay the dead letter
   - Show successful delivery

6. **Show Multi-Endpoint Fan-out** (30s)
   - One event → multiple endpoints
   - Independent timelines

**Total: ~5-6 minutes**

## Next Steps

- Explore the codebase structure
- Review retry logic in `internal/delivery/retry.go`
- Examine worker pool in `internal/delivery/worker.go`
- Check frontend real-time updates in `frontend/src/hooks/usePolling.ts`

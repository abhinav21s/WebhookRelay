# WebhookRelay - Setup Checklist

Use this checklist to verify your setup is complete and working.

## ✅ Prerequisites

- [ ] Go 1.22+ installed and in PATH
  ```bash
  go version  # Should show 1.22 or higher
  ```

- [ ] Node.js 18+ installed and in PATH
  ```bash
  node --version  # Should show v18 or higher
  npm --version
  ```

- [ ] Git installed
  ```bash
  git --version
  ```

- [ ] Supabase account created at https://supabase.com
- [ ] Code editor (VS Code recommended)

## ✅ Supabase Setup

- [ ] Created new Supabase project named "webhookrelay"
- [ ] Database provisioned (green status in dashboard)
- [ ] SQL tables created (endpoints, events, delivery_attempts)
- [ ] Indexes created successfully
- [ ] Copied Project URL from Settings → API
- [ ] Copied Anon Key from Settings → API
- [ ] Copied Database Connection String from Settings → Database

## ✅ Backend Setup

- [ ] `.env` file created from `.env.example`
- [ ] `SUPABASE_URL` configured in `.env`
- [ ] `SUPABASE_KEY` configured in `.env`
- [ ] `SUPABASE_DB_URL` configured in `.env`
- [ ] `go mod download` completed successfully
- [ ] No errors when running `go run cmd/api/main.go`

## ✅ Frontend Setup

- [ ] `frontend/.env` file created from `frontend/.env.example`
- [ ] `VITE_API_URL` set to `http://localhost:8080`
- [ ] `cd frontend && npm install` completed successfully
- [ ] No errors when running `npm run dev`

## ✅ Mock Receiver Setup

- [ ] `go run cmd/mock/main.go` starts without errors
- [ ] Listening on port 9090
- [ ] Default mode is "success"

## ✅ Connectivity Tests

- [ ] Backend API responds at http://localhost:8080/health
  ```bash
  curl http://localhost:8080/health
  # Expected: {"status":"ok"}
  ```

- [ ] Mock receiver responds at http://localhost:9090
  ```bash
  curl -X POST http://localhost:9090/webhook -H "Content-Type: application/json" -d '{"test":true}'
  # Should see logs in mock terminal
  ```

- [ ] Frontend loads at http://localhost:5173
- [ ] No console errors in browser DevTools

## ✅ Database Connection

- [ ] Backend logs show "Connected to Supabase successfully"
- [ ] Backend logs show "Worker pool started with 10 workers"
- [ ] No database connection errors

## ✅ Core Features Test

### Endpoint Management
- [ ] Can create new endpoint via UI
- [ ] Endpoint appears in list immediately
- [ ] Can view endpoint details
- [ ] Can toggle endpoint active/inactive
- [ ] Can delete endpoint

### Event Publishing
- [ ] Can send test event from endpoint detail page
- [ ] Event appears in Events list
- [ ] Event status updates in real-time
- [ ] Dashboard metrics update

### Delivery & Retry
- [ ] Test event delivers successfully (200)
- [ ] Timeline shows single successful attempt
- [ ] Latency recorded (~100-200ms)

### Failure Mode Testing
- [ ] Can change mock mode to "500"
- [ ] Event triggers retries (watch timeline)
- [ ] See 5 attempts with exponential backoff
- [ ] Event moves to Dead Letter after max attempts
- [ ] Delays visible: immediate, +1s, +2s, +4s, +8s

### Replay Feature
- [ ] Dead letter appears in Dead Letters page
- [ ] Can replay event
- [ ] New event created with same payload
- [ ] Successful delivery after replay

## ✅ UI/UX Checks

### Dashboard
- [ ] Metrics cards display correctly
- [ ] Success rate chart renders
- [ ] Recent events list populates
- [ ] Endpoint health indicators show
- [ ] Live polling updates data

### Endpoints Page
- [ ] Table displays all endpoints
- [ ] Status badges correct colors
- [ ] Event types displayed
- [ ] Actions buttons work
- [ ] Create modal functional

### Events Page
- [ ] Events list displays
- [ ] Status badges update in real-time
- [ ] Pagination/scrolling works
- [ ] Event details link works

### Event Detail
- [ ] Event ID and type displayed
- [ ] Payload formatted correctly (JSON)
- [ ] Delivery timeline renders
- [ ] Timeline items animate in
- [ ] Replay button functional

### Dead Letters Page
- [ ] Failed events listed
- [ ] Replay button works
- [ ] Inspect link navigates correctly

## ✅ API Endpoints

Test all endpoints with cURL:

- [ ] `POST /api/webhooks` - Create endpoint
- [ ] `GET /api/webhooks` - List endpoints
- [ ] `GET /api/webhooks/:id` - Get endpoint
- [ ] `PATCH /api/webhooks/:id` - Update endpoint
- [ ] `DELETE /api/webhooks/:id` - Delete endpoint
- [ ] `POST /api/events` - Publish event
- [ ] `GET /api/events` - List events
- [ ] `GET /api/events/:id` - Get event detail
- [ ] `POST /api/events/:id/replay` - Replay event
- [ ] `GET /api/deliveries?event_id=...` - List deliveries
- [ ] `GET /api/metrics` - Get metrics
- [ ] `POST /api/mock/control` - Set mock mode
- [ ] `GET /api/mock/mode` - Get mock mode

## ✅ Documentation

- [ ] README.md reviewed
- [ ] SETUP.md followed successfully
- [ ] RUNNING.md test scenarios work
- [ ] QUICK_START.md validated
- [ ] ARCHITECTURE.md understood
- [ ] PROJECT_SUMMARY.md reviewed

## ✅ Code Quality

- [ ] No compilation errors in Go code
- [ ] No TypeScript errors in frontend
- [ ] No console errors in browser
- [ ] All imports resolve correctly
- [ ] `.gitignore` configured properly

## ✅ Production Readiness (Optional)

- [ ] Environment variables secured
- [ ] Supabase RLS policies considered
- [ ] API rate limiting considered
- [ ] Authentication strategy planned
- [ ] Monitoring setup planned
- [ ] Logging strategy defined
- [ ] Backup strategy for Supabase
- [ ] Error tracking service integrated

## 🎯 Demo Preparation

For demonstrating to others:

- [ ] Clean database (delete test data)
- [ ] All services running
- [ ] Browser DevTools closed
- [ ] Terminal windows organized
- [ ] Mock set to "success" mode initially
- [ ] Test endpoints ready with localhost:9090
- [ ] Practiced 5-minute demo flow

## 🎯 5-Minute Demo Flow

1. [ ] Show Dashboard (30s)
   - Point out metrics
   - Show endpoint health

2. [ ] Create Endpoint (30s)
   - Name: "Demo Service"
   - URL: http://localhost:9090/webhook
   - Event: test.demo

3. [ ] Send Successful Event (1min)
   - Send test event
   - Navigate to Events
   - Click event to show timeline
   - Point out latency and success

4. [ ] Demonstrate Retry Logic (2min) ⭐
   - Set mock to "500" mode
   - Send test event
   - Watch timeline update live
   - Point out exponential backoff delays
   - Show Dead Letter status

5. [ ] Replay Event (1min)
   - Switch mock back to "success"
   - Navigate to Dead Letters
   - Replay event
   - Show successful delivery

## 🐛 Common Issues Resolved

- [ ] Port 8080 available (no conflicts)
- [ ] Port 9090 available (no conflicts)
- [ ] Port 5173 available (no conflicts)
- [ ] Firewall allows local connections
- [ ] Supabase project not paused
- [ ] Database password correct
- [ ] IP whitelisted in Supabase (if needed)
- [ ] CORS configured correctly
- [ ] Browser cache cleared if UI issues

## 📝 Notes

Write any custom configurations or observations:

```
_____________________________________________
_____________________________________________
_____________________________________________
_____________________________________________
```

---

**All checked? You're ready to go! 🚀**

If any items are unchecked, refer to the relevant documentation file for help.

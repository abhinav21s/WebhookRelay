# Fix Instructions - Endpoint Creation

## ✅ Changes Applied

I've updated `frontend/src/pages/Endpoints.tsx` with:
- ✅ Error handling and display
- ✅ Loading state ("Creating..." button)
- ✅ Console logging for debugging
- ✅ Validation for event types

## 🔧 Steps to Fix

### 1. Restart Frontend (REQUIRED)

Since you just created the `.env` file, you **MUST restart** the frontend:

```bash
# Go to Terminal 3 where frontend is running
# Press Ctrl+C to stop it

# Then restart:
cd frontend
npm run dev
```

**Why?** Vite doesn't reload when `.env` is created/changed. You must restart!

---

### 2. Clear Browser Cache

After restarting frontend:

1. Open browser (http://localhost:5173)
2. Press `Ctrl + Shift + R` (hard refresh)
3. Or press `F12` → Application tab → Clear storage → Clear site data

---

### 3. Test Backend Directly

Run this command to verify backend works:

```bash
test_endpoint_creation.bat
```

This will:
- Create an endpoint via API
- List all endpoints
- Show if backend is working

**If this works:** Backend is fine, issue is in frontend
**If this fails:** Backend/database issue

---

### 4. Open Browser Console

1. Open http://localhost:5173
2. Press `F12` to open DevTools
3. Go to **Console** tab
4. Keep it open

---

### 5. Try Creating Endpoint

1. Click **"Endpoints"** in navigation
2. Click **"New Endpoint"** button
3. Fill in:
   - Name: `Test Endpoint`
   - URL: `http://localhost:9090/webhook`
   - Event Types: `test.event`
4. Click **"Create"**

---

### 6. Check Console Output

You should see in console:

**Success case:**
```
Creating endpoint with data: {name: "Test Endpoint", url: "...", event_types: ["test.event"]}
✅ Endpoint created successfully
```

**Error case:**
```
Creating endpoint with data: {...}
❌ Create endpoint error: [error message]
```

---

### 7. Check Network Tab

1. In DevTools, go to **Network** tab
2. Clear all requests
3. Try creating endpoint again
4. Look for `webhooks` request

**Check:**
- Status code (should be 201)
- Response (should have endpoint data)
- If status is 500/400, click to see error details

---

## 🐛 Common Issues & Solutions

### Issue 1: "No endpoints appear after creation"

**Check if endpoint was actually created:**
```bash
curl http://localhost:8080/api/webhooks
```

**If endpoint exists in API but not in UI:**

1. Check browser console for React Query errors
2. Manually refresh page (F5)
3. Check if `refetchInterval` is working (should auto-refresh every 3 seconds)

**Solution:**
```typescript
// This is already in the code, but verify it's there:
const { data: endpoints = [] } = useQuery({
  queryKey: ['endpoints'],
  queryFn: fetchEndpoints,
  refetchInterval: 3000, // ← This line is important
})
```

---

### Issue 2: "Error: Failed to fetch"

**Cause:** Backend not running or wrong API URL

**Check:**
```bash
# Check backend
curl http://localhost:8080/health

# Check frontend .env
type frontend\.env
```

Should show: `VITE_API_URL=http://localhost:8080`

**Solution:**
1. Restart backend: `go run cmd/api/main.go`
2. Verify frontend .env has correct URL
3. Restart frontend

---

### Issue 3: "CORS Error"

**Full error in console:**
```
Access to fetch at 'http://localhost:8080/api/webhooks' 
from origin 'http://localhost:5173' has been blocked by CORS policy
```

**Solution:**
Backend already has CORS configured. Restart backend:
```bash
go run cmd/api/main.go
```

---

### Issue 4: "500 Internal Server Error"

**Cause:** Database connection issue

**Check Terminal 1 (backend) logs for:**
```
Failed to connect to database
```

**Solution:**
1. Verify Supabase project is not paused
2. Check `.env` has correct `SUPABASE_DB_URL`
3. Run the SQL schema: `supabase_schema.sql`

---

### Issue 5: Button says "Creating..." forever

**Cause:** Request hanging or never completing

**Solution:**
1. Check Network tab to see if request completed
2. Look for timeout errors
3. Restart backend and try again

---

## 🎯 Quick Fix Script

Run this to test everything:

```bash
# 1. Test backend API
test_endpoint_creation.bat

# 2. If backend works, restart frontend
cd frontend
# Ctrl+C to stop
npm run dev

# 3. Open browser with console
start http://localhost:5173
# Then press F12
```

---

## 🔍 Debug Checklist

Before asking for help, verify:

- [ ] Backend running (Terminal 1 shows "Worker pool started")
- [ ] Mock running (Terminal 2 shows "Mock receiver started")
- [ ] Frontend running (Terminal 3 shows "Local: http://localhost:5173")
- [ ] Backend health OK: `curl http://localhost:8080/health`
- [ ] Frontend .env exists with `VITE_API_URL=http://localhost:8080`
- [ ] Frontend restarted after creating .env
- [ ] Browser console open (F12)
- [ ] Network tab shows POST to `/api/webhooks`
- [ ] Manual API test works: `test_endpoint_creation.bat`

---

## ✨ Expected Behavior

**After clicking Create:**

1. Button changes to "Creating..."
2. Button becomes disabled
3. Console shows: "Creating endpoint with data: {...}"
4. Network tab shows POST to `/api/webhooks`
5. Response status: 201 Created
6. Console shows: "✅ Endpoint created successfully"
7. Modal closes
8. Endpoint appears in table (within 3 seconds due to auto-refresh)

**If endpoint doesn't appear immediately:**
- Wait 3 seconds (auto-refresh)
- Or press F5 to manually refresh

---

## 📞 Still Not Working?

Share these details:

1. **Browser console output** (after trying to create)
2. **Network tab screenshot** (the `/api/webhooks` request)
3. **Backend Terminal 1 output** (any errors?)
4. **Result of:** `curl http://localhost:8080/api/webhooks`

With these, I can help you fix the exact issue!

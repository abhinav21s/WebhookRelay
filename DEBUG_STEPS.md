# Debug Steps - Endpoints Not Showing

## Current Situation
- ✅ Database is connected
- ✅ Endpoints are created and stored in Supabase
- ✅ You can see them in Supabase table
- ❌ But API returns `[]` empty array
- ❌ Frontend shows no endpoints

## I've Added Debug Logging

### What I Changed:
1. Added logging to `CreateEndpoint` - shows when endpoint is created
2. Added logging to `ListEndpoints` - shows:
   - How many rows found in database
   - Any errors scanning rows
   - How many endpoints returned

## 🔧 Next Steps

### 1. Restart Backend

**IMPORTANT:** You must restart backend to get the new logging

In Terminal 1:
- Press `Ctrl+C` to stop
- Run: `go run cmd/api/main.go`
- Wait for: "Worker pool started with 10 workers"

### 2. Test GET Endpoints

```bash
curl http://localhost:8080/api/webhooks
```

### 3. Check Terminal 1 Logs

You should see something like:
```
Found 3 rows, returning 3 endpoints
```

**Or if there's an error:**
```
Error scanning row 1: sql: Scan error on column index 4...
```

### 4. Share the Output

**Tell me what you see in Terminal 1:**
- The log lines that appear
- Especially "Found X rows, returning Y endpoints"
- And any "Error scanning row" messages

---

## 🎯 Most Likely Issues

### Issue 1: Data Type Mismatch

**Symptom:** Logs show "Found 3 rows, returning 0 endpoints"

**Cause:** `event_types` column has wrong data type or format

**Solution:** 
Check how data is stored in Supabase:
```sql
SELECT id, name, event_types, pg_typeof(event_types) as type 
FROM endpoints 
LIMIT 1;
```

Should show type as `text[]` (array)

### Issue 2: NULL Values

**Symptom:** Error scanning row

**Cause:** Some column is NULL but Go expects a value

**Solution:**
Check for NULLs:
```sql
SELECT 
  id, 
  name IS NULL as name_null,
  url IS NULL as url_null,
  secret IS NULL as secret_null,
  event_types IS NULL as event_types_null,
  is_active IS NULL as is_active_null
FROM endpoints;
```

### Issue 3: Wrong Table

**Symptom:** Found 0 rows (but you know data exists)

**Cause:** Querying wrong database/schema

**Solution:**
Verify table has data:
```sql
SELECT COUNT(*) FROM endpoints;
```

Should show > 0

---

## 🔍 Detailed Debug Process

### Step 1: Check Data in Supabase

Go to Supabase → SQL Editor:

```sql
SELECT * FROM endpoints ORDER BY created_at DESC LIMIT 5;
```

**Copy the output and share it with me**

### Step 2: Check Backend Connection

In Terminal 1, you should see on startup:
```
Connected to Supabase successfully
```

### Step 3: Test API Call

```bash
curl http://localhost:8080/api/webhooks
```

Check Terminal 1 immediately after - should show:
```
Found X rows, returning Y endpoints
```

### Step 4: Create New Endpoint

Via browser:
1. Try creating endpoint
2. Check Terminal 1 logs

Should see:
```
Creating endpoint: name=Test, url=..., event_types=[...]
✅ Successfully created endpoint with ID: xxx-xxx-xxx
```

Then immediately:
```bash
curl http://localhost:8080/api/webhooks
```

Check Terminal 1 - should show the endpoint now!

---

## 🐛 Common Scan Errors

### Error: "Scan error on column index 4"

Column 4 is `event_types` (array type)

**Problem:** Data stored as string instead of array

**Fix:**
```sql
-- Check current data type
SELECT pg_typeof(event_types) FROM endpoints LIMIT 1;

-- If it's 'text' instead of 'text[]', fix the table:
ALTER TABLE endpoints 
ALTER COLUMN event_types TYPE text[] 
USING string_to_array(event_types, ',');
```

### Error: "Scan error: unsupported Scan"

**Problem:** Go type doesn't match PostgreSQL type

**Fix:** Check models.go has:
```go
EventTypes pq.StringArray `json:"event_types"`
```

Not:
```go
EventTypes []string `json:"event_types"`
```

---

## ✅ Expected Output After Fix

### Terminal 1 (Backend):
```
GET /api/webhooks - 200 (15ms)
Found 3 rows, returning 3 endpoints
```

### curl command:
```json
[
  {
    "id": "uuid-here",
    "name": "Test",
    "url": "http://localhost:9090/webhook",
    "event_types": ["test.event", "payment.completed"],
    "is_active": true,
    ...
  }
]
```

### Browser Console:
```
🔍 Fetching endpoints from: http://localhost:8080/api/webhooks
✅ Fetched endpoints: 3
```

### Browser UI:
- Table shows 3 endpoints
- Each with name, URL, event types visible

---

## 🚀 Quick Test Script

I'll create a test script for you...

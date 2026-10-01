# Quick Start Guide

Get WebhookRelay running in 5 minutes!

## Prerequisites Checklist

- [ ] Go 1.22+ installed
- [ ] Node.js 18+ installed
- [ ] Supabase account created
- [ ] Git installed

## Step-by-Step Setup

### 1. Create Supabase Project (5 minutes)

1. Go to https://app.supabase.com/
2. Click "New Project"
3. Set name: `webhookrelay`
4. Set database password (save it!)
5. Wait for provisioning (~2 min)

### 2. Create Database Tables

In Supabase SQL Editor, run:

```sql
-- Endpoints
CREATE TABLE endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    event_types TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Events
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Delivery attempts
CREATE TABLE delivery_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    endpoint_id UUID NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    attempt_number INT NOT NULL,
    status_code INT,
    latency_ms INT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_created_at ON events(created_at DESC);
CREATE INDEX idx_delivery_attempts_event_id ON delivery_attempts(event_id);
CREATE INDEX idx_delivery_attempts_endpoint_id ON delivery_attempts(endpoint_id);
CREATE INDEX idx_delivery_attempts_created_at ON delivery_attempts(created_at DESC);
```

### 3. Configure Environment

Copy `.env.example` to `.env`:

```bash
copy .env.example .env
```

Edit `.env` with your Supabase credentials:

```env
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_KEY=your-anon-key
SUPABASE_DB_URL=postgresql://postgres:password@db.xxxxx.supabase.co:5432/postgres
```

### 4. Install Dependencies

**Backend:**
```bash
go mod download
```

**Frontend:**
```bash
cd frontend
npm install
cd ..
```

### 5. Start Everything

**Option A - Automated (Windows):**
```bash
start.bat
```

**Option B - Manual:**

Terminal 1 (Backend):
```bash
go run cmd/api/main.go
```

Terminal 2 (Mock Receiver):
```bash
go run cmd/mock/main.go
```

Terminal 3 (Frontend):
```bash
cd frontend
npm run dev
```

### 6. Access Dashboard

Open browser: http://localhost:5173

## First Steps

### Create Your First Endpoint

1. Click **"Endpoints"** in navigation
2. Click **"New Endpoint"**
3. Fill in:
   - Name: `Test Endpoint`
   - URL: `http://localhost:9090/webhook`
   - Event Types: `test.event`
4. Click **"Create"**

### Send a Test Event

1. Click on the endpoint name
2. Set **Failure Mode** to `"Success (200)"`
3. Click **"Send Test Event"**
4. Click **"Send Event"**
5. Navigate to **Events** to see delivery status

### Watch Retry Logic (The Star Feature!)

1. Go back to endpoint detail
2. Change **Failure Mode** to `"500 Internal Server Error"`
3. Send another test event
4. Click on the event to see **Delivery Timeline**
5. Watch attempts appear with exponential backoff:
   - Attempt 1: Failed (500) - immediate
   - Attempt 2: Failed (500) - 1s later
   - Attempt 3: Failed (500) - 2s later
   - Attempt 4: Failed (500) - 4s later
   - Attempt 5: Failed (500) - 8s later
   - Status: Dead Letter

### Replay a Failed Event

1. Go to **"Dead Letters"**
2. Switch failure mode back to **"Success"**
3. Click **"Replay"** on the dead letter
4. Watch successful delivery!

## Troubleshooting

### "Failed to connect to database"
- Check Supabase project is not paused
- Verify database URL and password in `.env`
- Ensure IP is whitelisted in Supabase settings

### Frontend can't connect to API
- Verify backend is running on port 8080
- Check `frontend/.env` has correct API URL
- Check browser console for CORS errors

### Events not delivering
- Ensure endpoint is Active (green badge)
- Verify mock receiver is running on port 9090
- Check backend logs for errors

## API Examples

### Create Endpoint
```bash
curl -X POST http://localhost:8080/api/webhooks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Service",
    "url": "http://localhost:9090/webhook",
    "event_types": ["payment.completed"]
  }'
```

### Publish Event
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "payment.completed",
    "payload": {"amount": 99.99}
  }'
```

## Next Steps

- Read [RUNNING.md](./RUNNING.md) for detailed feature testing
- Review code structure in [README.md](./README.md)
- Explore retry logic in `internal/delivery/retry.go`
- Check dashboard real-time updates

## Need Help?

- Check logs in terminal windows
- Review Supabase dashboard for database issues
- Ensure all ports (8080, 9090, 5173) are available
- Verify Go and Node versions meet requirements

Happy webhook relaying! 🚀

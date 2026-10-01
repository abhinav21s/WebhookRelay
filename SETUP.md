# Setup Instructions

Complete guide to set up WebhookRelay on your local machine.

## Prerequisites

- **Go**: 1.22 or higher ([Download](https://golang.org/dl/))
- **Node.js**: 18 or higher ([Download](https://nodejs.org/))
- **Supabase Account**: Free tier is sufficient ([Sign up](https://supabase.com/))
- **Git**: For cloning the repository

## Step 1: Clone the Repository

```bash
git clone <repository-url>
cd webhookrelay
```

## Step 2: Set Up Supabase

### Create a Supabase Project

1. Go to [Supabase Dashboard](https://app.supabase.com/)
2. Click "New Project"
3. Choose organization and set:
   - **Name**: webhookrelay
   - **Database Password**: (save this!)
   - **Region**: Choose closest to you
4. Wait for project to provision (~2 minutes)

### Create Database Tables

1. In your Supabase project, go to **SQL Editor**
2. Run the following SQL:

```sql
-- Endpoints table
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

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Delivery attempts table
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

-- Indexes for performance
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_created_at ON events(created_at DESC);
CREATE INDEX idx_delivery_attempts_event_id ON delivery_attempts(event_id);
CREATE INDEX idx_delivery_attempts_endpoint_id ON delivery_attempts(endpoint_id);
CREATE INDEX idx_delivery_attempts_created_at ON delivery_attempts(created_at DESC);
```

### Get Your Supabase Credentials

1. In your Supabase project, go to **Project Settings** → **API**
2. Copy these values:
   - **Project URL**: e.g., `https://xxxxx.supabase.co`
   - **Project API Key (anon/public)**: The public key
3. Go to **Project Settings** → **Database**
4. Copy the **Connection String** (URI format)

## Step 3: Configure Environment Variables

### Backend Configuration

1. In the project root, create `.env`:

```bash
cp .env.example .env
```

2. Edit `.env` with your Supabase credentials:

```env
# Server Configuration
PORT=8080
ENV=development

# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
SUPABASE_DB_URL=postgresql://postgres:your-password@db.your-project.supabase.co:5432/postgres

# Delivery Engine
WORKER_POOL_SIZE=10
MAX_RETRY_ATTEMPTS=5
BASE_RETRY_DELAY=1s
MAX_RETRY_DELAY=30s
REQUEST_TIMEOUT=5s

# Mock Receiver (for testing)
MOCK_RECEIVER_PORT=9090
```

### Frontend Configuration

1. Create `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080
```

## Step 4: Install Dependencies

### Backend

```bash
go mod download
```

If `go.mod` doesn't exist yet, initialize:

```bash
go mod init github.com/yourusername/webhookrelay
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors
go get github.com/joho/godotenv
go get github.com/lib/pq
```

### Frontend

```bash
cd frontend
npm install
```

## Step 5: Verify Setup

### Test Supabase Connection

```bash
go run cmd/api/main.go
```

You should see:
```
2024/01/01 12:00:00 Starting WebhookRelay API on :8080
2024/01/01 12:00:00 Connected to Supabase successfully
2024/01/01 12:00:00 Worker pool started with 10 workers
```

### Test Frontend

```bash
cd frontend
npm run dev
```

You should see:
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

## Step 6: Initial Data (Optional)

Create a test endpoint via API:

```bash
curl -X POST http://localhost:8080/api/webhooks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Endpoint",
    "url": "http://localhost:9090/webhook",
    "event_types": ["test.event"],
    "secret": "test-secret-123"
  }'
```

## Troubleshooting

### Connection Issues

- Verify Supabase project is not paused (free tier auto-pauses after 1 week inactivity)
- Check if database URL includes correct password
- Ensure your IP is allowed (Supabase → Project Settings → Database → Connection Pooling)

### Port Already in Use

```bash
# Change PORT in .env to a different value
PORT=8081
```

### Frontend Can't Connect to API

- Verify backend is running on port 8080
- Check VITE_API_URL in frontend/.env
- Ensure no CORS issues (CORS is configured in backend)

## Next Steps

Proceed to [RUNNING.md](./RUNNING.md) to learn how to run and test all features.

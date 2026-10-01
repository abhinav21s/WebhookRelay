# WebhookRelay

A reliable, production-quality webhook delivery system with intelligent retry logic, real-time observability, and a beautiful dark-themed dashboard.

## What is WebhookRelay?

WebhookRelay is a lightweight webhook delivery platform that solves the core problem of reliably notifying external services. When endpoints go down, return errors, or time out, WebhookRelay handles queuing, concurrent delivery, exponential backoff retries, and provides full observability through a real-time dashboard.

## Architecture Overview

```
┌─────────────┐
│   React     │
│  Dashboard  │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────┐      ┌──────────────┐
│   Go API        │◄────►│  Supabase    │
│   (Chi Router)  │      │ (PostgreSQL) │
└────────┬────────┘      └──────────────┘
         │
         ▼
┌─────────────────┐
│  Worker Pool    │
│  - Goroutines   │
│  - Channel Queue│
│  - Exp. Backoff │
└────────┬────────┘
         │ HTTP POST
         ▼
┌─────────────────┐
│ Target Endpoints│
└─────────────────┘
```

## Key Features

- **Endpoint Management**: Register webhooks, subscribe to event types, enable/disable without data loss
- **Event Publishing**: Simple JSON API with automatic fan-out to matching endpoints
- **Reliable Delivery**: Worker pool with exponential backoff (1s → 2s → 4s → 8s → 30s, max 5 attempts)
- **HMAC Signing**: Every webhook signed with HMAC-SHA256 for verification
- **Full Observability**: Detailed delivery timeline, success rates, latency tracking
- **Dead Letter Queue**: Permanent failure tracking with one-click replay
- **Live Demo Mode**: Controllable mock receiver to demonstrate retry behavior
- **Real-time Dashboard**: Dark-themed UI with live status updates

## Tech Stack

**Backend:**
- Go 1.22+
- Chi HTTP router
- Supabase (PostgreSQL)
- Worker pool with buffered channels
- HMAC-SHA256 signing

**Frontend:**
- React 18 + Vite
- TypeScript
- Tailwind CSS
- TanStack Query (React Query)
- Recharts
- React Router

## Project Structure

```
webhookrelay/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Environment configuration
│   ├── database/
│   │   └── supabase.go             # Supabase client setup
│   ├── models/
│   │   └── models.go               # Data models
│   ├── handlers/
│   │   ├── webhooks.go             # Webhook endpoint handlers
│   │   ├── events.go               # Event handlers
│   │   ├── deliveries.go           # Delivery handlers
│   │   ├── metrics.go              # Metrics handlers
│   │   └── mock.go                 # Mock receiver handlers
│   ├── delivery/
│   │   ├── engine.go               # Delivery engine
│   │   ├── worker.go               # Worker pool
│   │   ├── retry.go                # Retry logic with exp. backoff
│   │   └── signer.go               # HMAC signing
│   └── middleware/
│       └── logging.go              # Request logging
├── frontend/
│   ├── src/
│   │   ├── components/             # React components
│   │   ├── pages/                  # Page components
│   │   ├── hooks/                  # Custom hooks
│   │   ├── services/               # API services
│   │   └── App.tsx                 # Main app
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml
├── .env.example
├── README.md                        # This file
├── SETUP.md                         # Setup instructions
└── RUNNING.md                       # How to run and test features
```

## Core Workflows

### 1. Register an Endpoint
```bash
POST /api/webhooks
{
  "name": "Payment Service",
  "url": "https://myservice.com/webhook",
  "event_types": ["payment.completed", "payment.failed"],
  "secret": "your-secret-key"
}
```

### 2. Publish an Event
```bash
POST /api/events
{
  "event_type": "payment.completed",
  "payload": {"order_id": "12345", "amount": 99.99}
}
```

### 3. Automatic Delivery
- System finds matching active endpoints
- Queues delivery jobs
- Worker pool processes with retries
- All attempts logged to Supabase

### 4. Monitor in Dashboard
- View real-time delivery status
- Inspect detailed timeline
- Replay failed events
- Test with controllable failures

## Quick Start

See [SETUP.md](./SETUP.md) for detailed setup instructions.
See [RUNNING.md](./RUNNING.md) for how to run and test all features.

## Security

- All outgoing webhooks signed with HMAC-SHA256
- Signature included in `X-Webhook-Signature` header
- Timestamp in `X-Webhook-Timestamp` to prevent replay attacks
- Endpoint secrets stored securely in Supabase

## Performance

- Configurable worker pool (default: 10 concurrent workers)
- Buffered channel queue for non-blocking event ingestion
- Exponential backoff prevents thundering herd
- Maximum retry delay: 30 seconds
- Maximum attempts: 5 per delivery

## License

MIT

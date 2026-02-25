# app-blacklist

Noona app that automatically declines marketplace bookings from customers in a Blacklist customer group.

## Tech Stack

- Go
- `labstack/echo/v4` - HTTP server
- `noona-hq/noona-sdk-go` - Noona platform SDK (OAuth, webhooks, customer groups)
- `go.mongodb.org/mongo-driver` - MongoDB for user/token storage
- `go.uber.org/zap` - structured logging
- `golang-jwt/jwt` - ID token verification
- Deployed via Helm / Docker

## Architecture / How it works

Standard Noona app pattern - OAuth-based install flow + webhook consumer.

**Install flow:**
1. Merchant installs via Noona App Store -> OAuth redirect to `/oauth/callback?code=...`
2. App exchanges code for token, fetches user, scaffolds Noona resources:
   - Creates a **Blacklist** customer group
   - Registers an **event creation webhook**
3. Stores user with OAuth tokens in MongoDB

**Webhook flow:**
- `POST /webhook` receives event creation callbacks from Noona
- Checks if the event was created via marketplace and if the attached customer is in the Blacklist customer group
- If both conditions met, automatically declines the appointment

**Uninstall:** `GET /oauth/callback?action=uninstall&id_token=...`

## Key interfaces / API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/oauth/callback` | OAuth install/uninstall handler |
| POST | `/webhook` | Noona event creation webhook receiver |
| GET | `/health` | Health check |

## Dependencies

- **Noona API** - customer groups, webhook registration, event management (via `noona-sdk-go`)
- MongoDB - user and token storage

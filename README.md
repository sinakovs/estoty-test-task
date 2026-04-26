# estoty-test-task

Nakama-based test project built with Go, Docker, and PostgreSQL.

Implemented RPC methods:

- `update_account_metadata`
- `get_game_config`
- `private_status`

## Additional Technologies And Libraries Used

- [Nakama](https://heroiclabs.com/nakama/) as the game server
- PostgreSQL as the database
- Docker and Docker Compose for local environment startup
- Go for the Nakama plugin implementation
- `github.com/heroiclabs/nakama-common` for Nakama runtime APIs
- `just` as an optional local task runner

## How To Start The Project

The project uses values from `.env`. A ready example is available in `.env.example`.

### Option 1: start with `just`

If `just` is installed:

```powershell
just nakama-up
```

Stop the project:

```powershell
just nakama-down
```

View logs:

```powershell
just nakama-logs
```

### Option 2: start without `just`

If `just` is not installed, run Docker Compose directly:

```powershell
docker compose -f docker-compose-postgres.yml -p nakama up --build -d
```

Stop the project:

```powershell
docker compose -f docker-compose-postgres.yml -p nakama down
```

View logs:

```powershell
docker compose -f docker-compose-postgres.yml -p nakama logs -f nakama
```

After startup, Nakama HTTP API is available at:

```text
http://127.0.0.1:7350
```

## How To Call The RPC Methods

### `get_game_config`

Returns the game configuration JSON.

```powershell
curl "http://127.0.0.1:7350/v2/rpc/get_game_config?http_key=defaulthttpkey&unwrap" `
  -X POST `
  -H "Content-Type: application/json" `
  -d "{}"
```

### `private_status`

Private RPC for server-to-server usage with `http_key`.

```powershell
curl "http://127.0.0.1:7350/v2/rpc/private_status?http_key=defaulthttpkey&unwrap" `
  -X POST `
  -H "Content-Type: application/json" `
  -d "{}"
```

### `update_account_metadata`

This RPC requires a user session token.

First, authenticate a user:

```powershell
curl "http://127.0.0.1:7350/v2/account/authenticate/device?create=true&username=nakama-test-user" `
  -X POST `
  -u "defaultkey:" `
  -H "Content-Type: application/json" `
  -d "{\"id\":\"nakama-test-device-001\"}"
```

Take the `token` value from the response, then call the RPC:

```powershell
curl "http://127.0.0.1:7350/v2/rpc/update_account_metadata?unwrap" `
  -X POST `
  -H "Authorization: Bearer <SESSION_TOKEN>" `
  -H "Content-Type: application/json" `
  -d "{\"metadata\":{\"favorite_color\":\"blue\",\"xp\":10,\"rarity\":\"rare\"}}"
```

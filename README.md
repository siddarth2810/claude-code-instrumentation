# Claude Code Instrumentation

Claude Code Instrumentation is a local tracing bridge for Claude Code.

It receives Claude Code hook events over HTTP.
It turns those events into OpenTelemetry traces.
It sends the traces to an OpenTelemetry Collector.
The collector sends them to Grafana Tempo.
You can inspect the traces in Grafana.

## What It Runs

The app has four local parts.

- `cmd/server.go` starts the Go HTTP server.
- The server listens on `http://localhost:10987`.
- `POST /hooks` receives Claude Code hook events.
- The server exports traces to `localhost:4318`.
- Docker Compose runs the collector, Tempo, and Grafana.

## Requirements

Install these first.

- Go
- Docker
- Docker Compose
- Claude Code (works with Ollama)

Check Go with this command.

```bash
go version
```

Check Docker with this command.

```bash
docker --version
```

Check Compose with this command.

```bash
docker compose version
```

## Run The App

### 1. Start The Trace Stack

Start the OpenTelemetry Collector, Tempo, and Grafana.

```bash
docker compose up -d
```

This starts these services.

- OpenTelemetry Collector on `localhost:4318`
- Tempo on `localhost:3200`
- Grafana on `localhost:3000`

Check that the containers are running.

```bash
docker compose ps
```

### 2. Start The Go Server

Run the hook server.

```bash
go run cmd/server.go
```

Keep this terminal open.
The server stops when you press `Ctrl+C`.

The server exposes these routes.

- `GET /` returns a small health response.
- `POST /hooks` receives hook events.

Check the health route from another terminal.

```bash
curl http://localhost:10987/
```

You should see this response.

```text
Hello, from server!
```

## Connect Claude Code

Claude Code must send hook events to this server.
Claude Code supports HTTP hooks.
HTTP hooks send the event JSON as a POST body.

Create or edit this file in the project where you run Claude Code.

```text
.claude/settings.local.json
```


Start a Claude Code session after this file is saved.
Use Claude Code normally.
Each supported hook event is sent to the Go server.
Use `/hooks` inside Claude Code to verify the loaded hooks.

## View Traces In Grafana

Open Grafana.

```text
http://localhost:3000
```

Use the default Grafana login if you have not changed it.

```text
username: admin
password: admin
```

Grafana may ask you to set a new password.

Add Tempo as a data source.

1. Open `Connections`.
2. Open `Data sources`.
3. Click `Add data source`.
4. Choose `Tempo`.
5. Set the URL to `http://tempo:3200`.
6. Click `Save & test`.

Open `Explore`.
Select the Tempo data source.
Search for traces from `claude-code-instrumentation`.

```
{.service.name="claude-code-instrumentation"}
```

## Stop The App

Stop the Go server with `Ctrl+C`.

Stop the Docker services with this command.

```bash
docker compose down
```

Remove the Tempo volume only if you want to delete saved traces.

```bash
docker compose down -v
```

## Run Tests

Run all tests with this command.

```bash
go test ./...
```


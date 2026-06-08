# Claude Code Instrumentation

Watch what Claude Code is doing — in real time — through visual traces in your browser.

---

## What It Does

Every time Claude Code takes an action (runs a tool, reads a file, etc.), this app captures that event and turns it into a visual timeline you can explore in Grafana.

**The flow:**
```
Claude Code → Go Server (port 10987) → OpenTelemetry Collector → Tempo → Grafana
```

---

## Before You Start

Install these four things:

| Tool | What it is |
|---|---|
| **Go** | Runs the hook server |
| **Docker + Docker Compose** | Runs the tracing stack |
| **Claude Code** | The AI coding tool being monitored |

---

## Setup (One Time)

### 1. Tell Claude Code to send events here

Create this file inside your project:

```
.claude/settings.json
```

Add the hook configuration pointing to `http://localhost:10987/hooks`. Start a new Claude Code session after saving — type `/hooks` to confirm it loaded.

### 2. Start the tracing stack

```bash
docker compose up -d
```

This starts three background services:
- **OpenTelemetry Collector** — receives trace data
- **Tempo** — stores the traces
- **Grafana** — lets you view them at [localhost:3000](http://localhost:3000)
  
<img width="1561" height="901" alt="screenshot-20260605-120232Z-selected" src="https://github.com/user-attachments/assets/db8ae015-bb27-4924-8b0a-706b4a1423af" />


### 3. Start the Go server

```bash
go run cmd/server.go
```

The server is now listening for Claude Code events at `localhost:10987`.

---

## Viewing Traces

1. Open **[localhost:3000](http://localhost:3000)** — log in with `admin` / `admin`
2. Go to **Connections → Data sources → Add data source**
3. Choose **Tempo**, set the URL to `http://tempo:3200`, and click **Save & test**
4. Open **Explore**, select the Tempo data source, and search:

```
{.service.name="claude-code-instrumentation"}
```

You'll see a timeline of everything Claude Code did.

---

## Stopping

```bash
# Stop everything
docker compose down

# Stop everything AND delete saved traces
docker compose down -v
```

Stop the Go server anytime with `Ctrl+C`.

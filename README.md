# YASB

Yet Another Service Bus: a small workflow engine written in Go. Workflows are stored in MongoDB, started by a trigger (e.g. an HTTP call) and run as a chain of actions.

> Early stage and mid-refactor: the engine isn't connected to the HTTP API yet. The API currently only stores and lists workflows.

## Requirements

- Go 1.27
- Docker (dev stack and the MongoDB-backed tests)

## Getting started

```bash
docker compose up -d --build
```

This starts the app with [air](https://github.com/air-verse/air) hot reload, plus MongoDB:

| Service | Address | Notes |
|---|---|---|
| API | `http://localhost:8080` | |
| MongoDB | `localhost:27017` | user/password `mongo`/`mongo`, data in `./mongo_data` |

Configuration comes from environment variables:

| Variable | Description |
|---|---|
| `MONGO_DB_URI` | MongoDB connection string |

## API

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/workflows` | List all workflows |
| `POST` | `/api/workflows` | Add a workflow (JSON body). `400` for an invalid body |

```bash
curl -X POST localhost:8080/api/workflows \
  -H 'Content-Type: application/json' \
  -d @data/wf.json
```

## Workflow format

A workflow has a `code`, a `name`, a `trigger` and a `startAction`. Actions are stored as an envelope, `{"type": ..., "action": {...}}`, where `type` selects the action implementation. Each action can point to the next one, forming a chain:

```json
{
  "code": "prv1",
  "name": "Prova 1",
  "trigger": { "type": "http", "endpoint": "/prova1", "verb": "POST", "schema": "{}" },
  "startAction": {
    "type": "http",
    "action": {
      "verb": "GET",
      "endpoint": "https://jsonplaceholder.typicode.com/posts",
      "nextAction": null
    }
  }
}
```

### Available actions

| Type | Fields | Description |
|---|---|---|
| `http` | `endpoint`, `verb`, `headers`, `body`, `nextAction` | Sends an HTTP request and reads a JSON response |

### Adding an action

1. Implement `models.Action` (`GetType`, `Exec`, `Next`, `Outputs`) in the `actions` package.
2. Add `json` and `bson` tags, and use `models.ActionEnvelope` for any field that holds another action.
3. Register it in `init()`: `models.RegisterAction("mytype", func() models.Action { return &MyAction{} })`.

## Project layout

| Path | Contents |
|---|---|
| `main.go` | Config, store connection, HTTP routes |
| `actions/` | Action implementations |
| `models/` | Workflow, `Action` interface, envelope and type registry |
| `engine/` | Runs an action chain; input message parsing |
| `store/` | MongoDB access |
| `handlers/` | Fiber HTTP handlers |
| `config/` | Environment-based configuration |
| `data/` | Sample workflow |

## Testing

```bash
make test      # readable output via gotestsum, never cached
make test-ci   # JUnit XML + coverage in reports/
```

Store and handler tests start a real MongoDB with [testcontainers](https://golang.testcontainers.org/). They need Docker and are skipped without it.

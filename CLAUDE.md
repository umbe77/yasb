# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

YASB ("Yet Another Service Bus") is an early-stage Go workflow engine: workflows are stored in MongoDB, triggered (e.g. by HTTP), and run as a chain of actions. It is mid-refactor, so many pieces are stubs or not yet wired together.

## Commands

- Run the full dev stack (app with air hot-reload + MongoDB): `docker compose up --build`. The repo is bind-mounted into the container at `/app`; air rebuilds to `./tmp/main` on `.go` changes. API on `:8080`, Mongo on `:27017` (user/pass `mongo`/`mongo`, data in `./mongo_data`).
- Build: `go build ./...`
- Tests: `make test` (gotestsum, uncached) or `make test-ci` (JUnit XML + coverage in `reports/`); single test: `go test ./engine -run TestParseMessage`. Store/handler tests start MongoDB via testcontainers (need Docker, skip otherwise). Avoid `./...` — the root-owned `mongo_data/` dir makes it print permission errors.
- Tests use plain `testing` (no assertion library); fakes are small structs in the test files. DB-backed tests (`store/`, `handlers/`) share one container per package via `TestMain` and skip when Docker is unavailable.
- Go 1.27, Fiber **v3** (handlers take `fiber.Ctx` by value, not `*fiber.Ctx`).

## Architecture

- `main.go` — loads config, connects the store, registers Fiber routes (`GET/POST /api/workflows`).
- `config/` — env-based config via `envconfig`, loaded once (`sync.Once`). `MONGO_DB_URI` is the only setting.
- `store/` — MongoDB access (database `yasb`, collection `workflows` via the `workflowsColl` const).
- `handlers/` — Fiber handlers built as closures over `*store.Storage`; errors returned as `models.ErrorMessage` (400 for invalid input, 500 for store failures).
- `models/` — `Workflow` (code, name, `Trigger`, `StartAction`), and the `Action` interface (`Exec`, `Next`, `Outputs`, `GetType`). Actions form a linked list via `Next()`.
  - Polymorphic actions: any field holding an action is a `models.ActionEnvelope` (wrapper with `.Action`; build with `models.Wrap(a)`). It implements JSON and BSON (un)marshalling as `{"type": "<GetType()>", "action": {...}}`; nil is JSON `null` / empty BSON doc.
  - Decoding resolves `type` through the registry in `models/action_registry.go` (`RegisterAction`/`NewAction`). `models` can't import `actions` (import cycle), so each action registers itself in its package `init()`, and `main.go` blank-imports `actions`. Anything decoding workflows (incl. tests) must import `actions`.
  - `data/wf.json` is a sample workflow in the envelope shape.
- `actions/` — concrete `Action` implementations, keyed by `GetType()` (currently only `http`: `HttpRequest`). New action: implement `models.Action`, add json+bson tags, use `models.ActionEnvelope` for next/child actions, and call `models.RegisterAction` in `init()`.
- `engine/` — `Engine.Execute` walks the action chain from `StartAction.Action` calling `Exec` with a shared `*models.Message`; `Init`/`AddWorflow` keep workflows by code (duplicates rejected). The engine isn't used by `main.go` yet. `engine/message.go` has a separate `engine.Message`/`TaskMessage` type and `ParseInput`, which recursively normalizes input JSON (whole floats → int, RFC3339 strings → `time.Time`) using the generic `stack` package.
- `engine/tasks/` — older `Task` interface based on `engine.Message`, overlapping with `models.Action`; part of the pending refactor.

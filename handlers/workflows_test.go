package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/umbe77/yasb/actions"
	"github.com/umbe77/yasb/handlers"
	"github.com/umbe77/yasb/models"
	"github.com/umbe77/yasb/store"
)

// storage is nil when Docker isn't available; tests then skip.
var (
	storage  *store.Storage
	startErr error
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := mongodb.Run(ctx, "mongo:noble")
	if err == nil {
		var uri string
		if uri, err = container.ConnectionString(ctx); err == nil {
			storage, err = store.NewStore(uri)
		}
	}
	startErr = err

	code := m.Run()

	if storage != nil {
		_ = storage.Client.Disconnect(ctx)
	}
	if container != nil {
		_ = testcontainers.TerminateContainer(container)
	}
	os.Exit(code)
}

func newApp(t *testing.T) *fiber.App {
	t.Helper()
	if storage == nil {
		t.Skipf("mongodb container unavailable: %v", startErr)
	}
	app := fiber.New()
	app.Get("/api/workflows", handlers.GetWorflows(storage))
	app.Post("/api/workflows", handlers.AddWorkflow(storage))
	return app
}

func do(t *testing.T, app *fiber.App, method, body string) (int, []byte) {
	t.Helper()
	req := httptest.NewRequest(method, "/api/workflows", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b
}

func TestAddThenGetWorkflow(t *testing.T) {
	app := newApp(t)
	raw, err := os.ReadFile("../data/wf.json")
	if err != nil {
		t.Fatal(err)
	}

	if status, body := do(t, app, http.MethodPost, string(raw)); status != http.StatusOK {
		t.Fatalf("POST status = %d, body = %s", status, body)
	}

	status, body := do(t, app, http.MethodGet, "")
	if status != http.StatusOK {
		t.Fatalf("GET status = %d, body = %s", status, body)
	}
	var wfs []models.Workflow
	if err := json.Unmarshal(body, &wfs); err != nil {
		t.Fatalf("decode %s: %v", body, err)
	}
	for _, wf := range wfs {
		if wf.Code != "prv1" {
			continue
		}
		h, ok := wf.StartAction.Action.(*actions.HttpRequest)
		if !ok || h.Verb != "GET" {
			t.Fatalf("start action: %#v", wf.StartAction.Action)
		}
		return
	}
	t.Fatalf("workflow prv1 not returned: %s", body)
}

func TestAddWorkflowBadRequest(t *testing.T) {
	app := newApp(t)
	cases := map[string]string{
		"empty body":   "",
		"invalid json": `{"code":`,
		"unknown type": `{"code":"x","startAction":{"type":"nope","action":{}}}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			status, resp := do(t, app, http.MethodPost, body)
			if status != http.StatusBadRequest {
				t.Fatalf("status = %d, body = %s", status, resp)
			}
			var msg models.ErrorMessage
			if err := json.Unmarshal(resp, &msg); err != nil || msg.Error == "" {
				t.Fatalf("error body = %s (%v)", resp, err)
			}
		})
	}
}

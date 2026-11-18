package store_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/umbe77/yasb/actions"
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
		if err := testcontainers.TerminateContainer(container); err != nil {
			fmt.Fprintln(os.Stderr, "terminate mongo container:", err)
		}
	}
	os.Exit(code)
}

func requireStore(t *testing.T) *store.Storage {
	t.Helper()
	if storage == nil {
		t.Skipf("mongodb container unavailable: %v", startErr)
	}
	return storage
}

func getByCode(t *testing.T, s *store.Storage, code string) models.Workflow {
	t.Helper()
	wfs, err := s.GetWorflows(context.Background(), fmt.Sprintf(`{"code":%q}`, code))
	if err != nil {
		t.Fatal(err)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow with code %q, got %d", code, len(wfs))
	}
	return wfs[0]
}

func TestAddAndGetWorkflow(t *testing.T) {
	s := requireStore(t)
	second := &actions.HttpRequest{Endpoint: "http://b", Verb: "POST", Body: map[string]any{"k": "v"}}
	first := &actions.HttpRequest{Endpoint: "http://a", Verb: "GET", NextAction: models.Wrap(second)}
	wf := models.NewWorkflow("chain", "Chain", models.Trigger{Type: "http", Endpoint: "/chain", Verb: "POST"}, first)

	if err := s.AddWorkflow(context.Background(), *wf); err != nil {
		t.Fatal(err)
	}

	all, err := s.GetWorflows(context.Background(), "{}")
	if err != nil || len(all) == 0 {
		t.Fatalf("GetWorflows({}): err=%v len=%d", err, len(all))
	}

	got := getByCode(t, s, "chain")
	if got.Name != "Chain" || got.Trigger.Endpoint != "/chain" {
		t.Fatalf("unexpected workflow: %+v", got)
	}
	a, ok := got.StartAction.Action.(*actions.HttpRequest)
	if !ok || a.Endpoint != "http://a" || a.Verb != "GET" {
		t.Fatalf("start action: %#v", got.StartAction.Action)
	}
	b, ok := a.Next().(*actions.HttpRequest)
	if !ok || b.Endpoint != "http://b" || b.Body["k"] != "v" {
		t.Fatalf("next action: %#v", a.Next())
	}
	if b.Next() != nil {
		t.Fatalf("expected end of chain, got %#v", b.Next())
	}
}

func TestLoadSampleWorkflow(t *testing.T) {
	s := requireStore(t)
	raw, err := os.ReadFile("../data/wf.json")
	if err != nil {
		t.Fatal(err)
	}
	var wf models.Workflow
	if err := json.Unmarshal(raw, &wf); err != nil {
		t.Fatal(err)
	}
	if err := s.AddWorkflow(context.Background(), wf); err != nil {
		t.Fatal(err)
	}

	got := getByCode(t, s, wf.Code)
	gotJSON, _ := json.Marshal(got)
	wantJSON, _ := json.Marshal(wf)
	if string(gotJSON) != string(wantJSON) {
		t.Fatalf("round trip mismatch:\n got  %s\n want %s", gotJSON, wantJSON)
	}
}

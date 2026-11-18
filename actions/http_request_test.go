package actions

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/umbe77/yasb/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func sampleWorkflow() *models.Workflow {
	second := &HttpRequest{Endpoint: "http://b", Verb: "POST", Headers: map[string]string{"X": "1"}}
	first := &HttpRequest{Endpoint: "http://a", Verb: "GET", NextAction: models.Wrap(second)}
	return models.NewWorkflow("c1", "n1", models.Trigger{Type: "http", Endpoint: "/x", Verb: "POST"}, first)
}

func assertChain(t *testing.T, wf models.Workflow) {
	t.Helper()
	first, ok := wf.StartAction.Action.(*HttpRequest)
	if !ok {
		t.Fatalf("start action: got %T", wf.StartAction.Action)
	}
	if first.Endpoint != "http://a" || first.Verb != "GET" {
		t.Fatalf("first: %+v", first)
	}
	second, ok := first.Next().(*HttpRequest)
	if !ok {
		t.Fatalf("next action: got %T", first.Next())
	}
	if second.Endpoint != "http://b" || second.Headers["X"] != "1" {
		t.Fatalf("second: %+v", second)
	}
	if second.Next() != nil {
		t.Fatalf("expected end of chain, got %T", second.Next())
	}
}

func TestWorkflowJSONRoundTrip(t *testing.T) {
	b, err := json.Marshal(sampleWorkflow())
	if err != nil {
		t.Fatal(err)
	}
	var wf models.Workflow
	if err := json.Unmarshal(b, &wf); err != nil {
		t.Fatal(err)
	}
	assertChain(t, wf)
}

func TestWorkflowBSONRoundTrip(t *testing.T) {
	b, err := bson.Marshal(sampleWorkflow())
	if err != nil {
		t.Fatal(err)
	}
	var wf models.Workflow
	if err := bson.Unmarshal(b, &wf); err != nil {
		t.Fatal(err)
	}
	assertChain(t, wf)
}

func TestSampleWorkflowFile(t *testing.T) {
	b, err := os.ReadFile("../data/wf.json")
	if err != nil {
		t.Fatal(err)
	}
	var wf models.Workflow
	if err := json.Unmarshal(b, &wf); err != nil {
		t.Fatal(err)
	}
	h, ok := wf.StartAction.Action.(*HttpRequest)
	if !ok || h.Verb != "GET" || h.Next() != nil {
		t.Fatalf("unexpected start action: %#v", wf.StartAction.Action)
	}
}

func TestEnvelopeEdgeCases(t *testing.T) {
	var e models.ActionEnvelope
	if err := json.Unmarshal([]byte(`{"type":"nope","action":{}}`), &e); err == nil {
		t.Fatal("expected error for unknown type")
	}
	if err := json.Unmarshal([]byte(`null`), &e); err != nil || e.Action != nil {
		t.Fatalf("null: err=%v action=%v", err, e.Action)
	}
	b, _ := bson.Marshal(models.ActionEnvelope{})
	if err := bson.Unmarshal(b, &e); err != nil || e.Action != nil {
		t.Fatalf("empty bson: err=%v action=%v", err, e.Action)
	}
}

func TestHttpRequestExec(t *testing.T) {
	var gotMethod, gotHeader string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotHeader = r.Header.Get("X-Test")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"n":2}`))
	}))
	defer srv.Close()

	h := &HttpRequest{
		Endpoint: srv.URL,
		Verb:     http.MethodPost,
		Headers:  map[string]string{"X-Test": "yes"},
		Body:     map[string]any{"k": "v"},
	}
	if err := h.Exec(context.Background(), &models.Message{}); err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost || gotHeader != "yes" || gotBody["k"] != "v" {
		t.Fatalf("server saw method=%s header=%s body=%v", gotMethod, gotHeader, gotBody)
	}
	if h.resp["ok"] != true || h.resp["n"] != float64(2) {
		t.Fatalf("resp = %v", h.resp)
	}
}

func TestHttpRequestExecErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	cases := map[string]*HttpRequest{
		"non-json response": {Endpoint: srv.URL, Verb: http.MethodGet},
		"unreachable":       {Endpoint: "http://127.0.0.1:1", Verb: http.MethodGet},
		"invalid verb":      {Endpoint: srv.URL, Verb: "BAD VERB"},
	}
	for name, h := range cases {
		t.Run(name, func(t *testing.T) {
			if err := h.Exec(context.Background(), &models.Message{}); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestHttpRequestOutputsAndType(t *testing.T) {
	last := &HttpRequest{}
	if out := last.Outputs(); len(out) != 0 {
		t.Fatalf("outputs without next: %v", out)
	}
	first := &HttpRequest{NextAction: models.Wrap(last)}
	if out := first.Outputs(); len(out) != 1 || out[0] != last {
		t.Fatalf("outputs with next: %v", out)
	}
	if first.GetType() != "http" {
		t.Fatalf("type = %q", first.GetType())
	}
}

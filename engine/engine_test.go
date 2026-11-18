package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/umbe77/yasb/models"
)

type recAction struct {
	name string
	err  error
	next models.Action
}

func (r *recAction) GetType() string { return "rec" }
func (r *recAction) Exec(ctx context.Context, m *models.Message) error {
	(*m)["order"] = append((*m)["order"].([]string), r.name)
	return r.err
}
func (r *recAction) Next() models.Action      { return r.next }
func (r *recAction) Outputs() []models.Action { return nil }

func newMsg() *models.Message { return &models.Message{"order": []string{}} }

func TestExecuteRunsChainInOrder(t *testing.T) {
	c := &recAction{name: "c"}
	b := &recAction{name: "b", next: c}
	a := &recAction{name: "a", next: b}
	wf := models.NewWorkflow("x", "x", models.Trigger{}, a)

	m := newMsg()
	if err := NewEngine().Execute(context.Background(), wf, m); err != nil {
		t.Fatal(err)
	}
	got := (*m)["order"].([]string)
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("order = %v", got)
	}
}

func TestExecuteStopsOnError(t *testing.T) {
	boom := errors.New("boom")
	c := &recAction{name: "c"}
	b := &recAction{name: "b", err: boom, next: c}
	a := &recAction{name: "a", next: b}
	wf := models.NewWorkflow("x", "x", models.Trigger{}, a)

	m := newMsg()
	if err := NewEngine().Execute(context.Background(), wf, m); !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if got := (*m)["order"].([]string); len(got) != 2 {
		t.Fatalf("order = %v, expected to stop after b", got)
	}
}

func TestExecuteNilStartAction(t *testing.T) {
	wf := models.NewWorkflow("x", "x", models.Trigger{}, nil)
	if err := NewEngine().Execute(context.Background(), wf, newMsg()); err != nil {
		t.Fatal(err)
	}
}

func TestInitAndAddWorkflow(t *testing.T) {
	e := NewEngine()
	e.Init([]models.Workflow{{Code: "a"}, {Code: "b"}})
	if len(e.workflows) != 2 {
		t.Fatalf("workflows = %d, want 2", len(e.workflows))
	}
	if err := e.AddWorflow(models.Workflow{Code: "c"}); err != nil {
		t.Fatalf("add new: %v", err)
	}
	if err := e.AddWorflow(models.Workflow{Code: "a"}); err == nil {
		t.Fatal("expected error adding duplicate")
	}
	if len(e.workflows) != 3 {
		t.Fatalf("workflows = %d, want 3", len(e.workflows))
	}
}

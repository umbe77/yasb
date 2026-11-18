package models

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeAction struct {
	Value      string         `json:"value" bson:"value"`
	NextAction ActionEnvelope `json:"nextAction" bson:"nextAction"`
}

func (f *fakeAction) GetType() string                            { return "fake" }
func (f *fakeAction) Exec(ctx context.Context, m *Message) error { return nil }
func (f *fakeAction) Next() Action                               { return f.NextAction.Action }
func (f *fakeAction) Outputs() []Action                          { return nil }

func init() {
	RegisterAction("fake", func() Action { return &fakeAction{} })
}

func TestNewAction(t *testing.T) {
	if _, err := NewAction("missing"); err == nil {
		t.Fatal("expected error for unknown type")
	}
	a, err := NewAction("fake")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := NewAction("fake")
	if _, ok := a.(*fakeAction); !ok {
		t.Fatalf("got %T", a)
	}
	if a == b {
		t.Fatal("factory should return a new instance each call")
	}
}

func TestRegisterActionDuplicatePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate registration")
		}
	}()
	RegisterAction("fake", func() Action { return &fakeAction{} })
}

func chain() ActionEnvelope {
	return Wrap(&fakeAction{Value: "a", NextAction: Wrap(&fakeAction{Value: "b"})})
}

func assertFakeChain(t *testing.T, e ActionEnvelope) {
	t.Helper()
	a, ok := e.Get().(*fakeAction)
	if !ok || a.Value != "a" {
		t.Fatalf("first: %#v", e.Get())
	}
	b, ok := a.Next().(*fakeAction)
	if !ok || b.Value != "b" || b.Next() != nil {
		t.Fatalf("second: %#v", a.Next())
	}
}

func TestEnvelopeJSON(t *testing.T) {
	data, err := json.Marshal(chain())
	if err != nil {
		t.Fatal(err)
	}
	var e ActionEnvelope
	if err := json.Unmarshal(data, &e); err != nil {
		t.Fatal(err)
	}
	assertFakeChain(t, e)

	nilData, _ := json.Marshal(ActionEnvelope{})
	if string(nilData) != "null" {
		t.Fatalf("nil action: got %s", nilData)
	}
}

func TestEnvelopeJSONEdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr bool
		wantNil bool
	}{
		{"null", `null`, false, true},
		{"empty type", `{"type":""}`, false, true},
		{"missing action", `{"type":"fake"}`, false, false},
		{"unknown type", `{"type":"nope","action":{}}`, true, true},
		{"malformed", `{"type":`, true, true},
		{"action wrong shape", `{"type":"fake","action":[1]}`, true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := Wrap(&fakeAction{Value: "stale"})
			err := json.Unmarshal([]byte(c.in), &e)
			if (err != nil) != c.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, c.wantErr)
			}
			if err == nil && (e.Action == nil) != c.wantNil {
				t.Fatalf("action = %#v, wantNil %v", e.Action, c.wantNil)
			}
			if err == nil && !c.wantNil {
				if f := e.Action.(*fakeAction); f.Value != "" {
					t.Fatalf("expected zero-valued action, got %#v", f)
				}
			}
		})
	}
}

func TestEnvelopeBSON(t *testing.T) {
	data, err := bson.Marshal(chain())
	if err != nil {
		t.Fatal(err)
	}
	var e ActionEnvelope
	if err := bson.Unmarshal(data, &e); err != nil {
		t.Fatal(err)
	}
	assertFakeChain(t, e)

	empty, _ := bson.Marshal(ActionEnvelope{})
	e = chain()
	if err := bson.Unmarshal(empty, &e); err != nil || e.Action != nil {
		t.Fatalf("empty doc: err=%v action=%#v", err, e.Action)
	}

	unknown, _ := bson.Marshal(bson.D{{Key: "type", Value: "nope"}})
	if err := bson.Unmarshal(unknown, &e); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestGetErrorMessage(t *testing.T) {
	m := GetErrorMessage("msg", errors.New("boom"))
	if m.Message != "msg" || m.Error != "boom" {
		t.Fatalf("got %+v", m)
	}
}

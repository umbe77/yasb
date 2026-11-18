package models

import (
	"bytes"
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ActionEnvelope wraps a polymorphic Action so it can be (de)serialized
// to/from JSON and BSON as {"type": "<GetType()>", "action": {...}}.
// A nil Action is encoded as JSON null and as an empty BSON document.
type ActionEnvelope struct {
	Action Action
}

func Wrap(a Action) ActionEnvelope { return ActionEnvelope{Action: a} }

func (e ActionEnvelope) Get() Action { return e.Action }

type jsonEnvelope struct {
	Type   string          `json:"type"`
	Action json.RawMessage `json:"action"`
}

type bsonEnvelope struct {
	Type   string        `bson:"type"`
	Action bson.RawValue `bson:"action"`
}

func (e ActionEnvelope) MarshalJSON() ([]byte, error) {
	if e.Action == nil {
		return []byte("null"), nil
	}
	raw, err := json.Marshal(e.Action)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonEnvelope{Type: e.Action.GetType(), Action: raw})
}

func (e *ActionEnvelope) UnmarshalJSON(data []byte) error {
	e.Action = nil
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var env jsonEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	if env.Type == "" {
		return nil
	}
	a, err := NewAction(env.Type)
	if err != nil {
		return err
	}
	if len(env.Action) > 0 {
		if err := json.Unmarshal(env.Action, a); err != nil {
			return err
		}
	}
	e.Action = a
	return nil
}

func (e ActionEnvelope) MarshalBSON() ([]byte, error) {
	if e.Action == nil {
		return bson.Marshal(bson.D{})
	}
	return bson.Marshal(bson.D{
		{Key: "type", Value: e.Action.GetType()},
		{Key: "action", Value: e.Action},
	})
}

func (e *ActionEnvelope) UnmarshalBSON(data []byte) error {
	e.Action = nil
	if len(data) == 0 {
		return nil
	}
	var env bsonEnvelope
	if err := bson.Unmarshal(data, &env); err != nil {
		return err
	}
	if env.Type == "" {
		return nil
	}
	a, err := NewAction(env.Type)
	if err != nil {
		return err
	}
	if env.Action.Type == bson.TypeEmbeddedDocument {
		if err := env.Action.Unmarshal(a); err != nil {
			return err
		}
	}
	e.Action = a
	return nil
}

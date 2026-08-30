package models

import "encoding/json"

type ActionEnvelope struct {
	Type   string          `json:"type"`
	Action json.RawMessage `json:"action"`
}

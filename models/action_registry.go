package models

import "fmt"

// ActionFactory builds an empty concrete Action ready to be decoded into.
type ActionFactory func() Action

var actionRegistry = map[string]ActionFactory{}

// RegisterAction associates an action type (as returned by GetType) with a factory.
// It is meant to be called from init() of packages implementing actions.
func RegisterAction(typ string, f ActionFactory) {
	if _, ok := actionRegistry[typ]; ok {
		panic(fmt.Sprintf("action type %q already registered", typ))
	}
	actionRegistry[typ] = f
}

// NewAction returns a new empty Action for the given type.
func NewAction(typ string) (Action, error) {
	f, ok := actionRegistry[typ]
	if !ok {
		return nil, fmt.Errorf("unknown action type %q", typ)
	}
	return f(), nil
}

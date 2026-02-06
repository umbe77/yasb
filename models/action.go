package models

import "context"

type Action interface {
	GetType() string
	Exec(ctx context.Context, m *Message) error
	Next() Action
	Outputs() []Action
}

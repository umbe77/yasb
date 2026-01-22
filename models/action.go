package models

import "context"

type Action interface {
	Exec(ctx context.Context, m *Message) error
	Next() Action
	Outputs() []Action
}
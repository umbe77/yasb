package tasks

import "github.com/umbe77/yasb/engine"

type Task interface {
	GetName() string
	GetDescription() string
	Execute(in engine.Message) (engine.Message, error)
}

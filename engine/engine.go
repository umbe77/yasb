package engine

import (
	"context"

	"github.com/umbe77/yasb/models"
)

type Engine struct {
	ctx context.Context
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Execute(ctx context.Context, wf *models.Workflow, m *models.Message) error {
	wf.Trigger.Exec(ctx, m)
	for currentAction := wf.Trigger.Next(); currentAction != nil; currentAction = currentAction.Next() {
		if err := currentAction.Exec(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
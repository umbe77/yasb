package engine

import (
	"context"
	"fmt"

	"github.com/umbe77/yasb/models"
)

type Engine struct {
	workflows map[string]models.Workflow
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Init(workflows []models.Workflow) {
	for _, wf := range workflows {
		e.workflows[wf.Code] = wf
	}
}

func (e *Engine) AddWorflow(wf models.Workflow) error {
	if _, ok := e.workflows[wf.Code]; !ok {
		e.workflows[wf.Code] = wf
	}
	return fmt.Errorf("Cannot add workflow %s, it is just added", wf.Code)
}

func (e *Engine) Execute(ctx context.Context, wf *models.Workflow, m *models.Message) error {
	for currentAction := wf.StartAction; currentAction != nil; currentAction = currentAction.Next() {
		if err := currentAction.Exec(ctx, m); err != nil {
			return err
		}
	}
	return nil
}


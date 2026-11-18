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
	return &Engine{workflows: make(map[string]models.Workflow)}
}

func (e *Engine) Init(workflows []models.Workflow) {
	for _, wf := range workflows {
		e.workflows[wf.Code] = wf
	}
}

func (e *Engine) AddWorflow(wf models.Workflow) error {
	if _, ok := e.workflows[wf.Code]; ok {
		return fmt.Errorf("cannot add workflow %s, it is already added", wf.Code)
	}
	e.workflows[wf.Code] = wf
	return nil
}

func (e *Engine) Execute(ctx context.Context, wf *models.Workflow, m *models.Message) error {
	for currentAction := wf.StartAction.Action; currentAction != nil; currentAction = currentAction.Next() {
		if err := currentAction.Exec(ctx, m); err != nil {
			return err
		}
	}
	return nil
}


package models

type Workflow struct {
	Code        string
	Name        string
	Trigger     Trigger
	StartAction Action
}

func NewWorkflow(code string, name string, trigger Trigger, act Action) *Workflow {
	return &Workflow{
		Code:        code,
		Name:        name,
		Trigger:     trigger,
		StartAction: act,
	}
}

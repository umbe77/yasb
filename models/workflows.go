package models

type Workflow struct {
	Trigger Action
	Actions []Action
}

func NewWorkflow(acts []Action) *Workflow {
	return &Workflow{
		Actions: acts,
	}
}

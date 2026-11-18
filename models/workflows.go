package models

type Workflow struct {
	Code        string         `json:"code" bson:"code"`
	Name        string         `json:"name" bson:"name"`
	Trigger     Trigger        `json:"trigger" bson:"trigger"`
	StartAction ActionEnvelope `json:"startAction" bson:"startAction"`
}

func NewWorkflow(code string, name string, trigger Trigger, act Action) *Workflow {
	return &Workflow{
		Code:        code,
		Name:        name,
		Trigger:     trigger,
		StartAction: Wrap(act),
	}
}

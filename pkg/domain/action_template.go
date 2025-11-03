package domain

type ActionArgs struct {
	VariableType string
	VariableName string
}

type ActionTemplate struct {
	Action     func()
	Parameters []ActionArgs
}

func NewActionTemplate(action func(), parameters []ActionArgs) *ActionTemplate {
	return &ActionTemplate{
		Action:     action,
		Parameters: parameters,
	}
}

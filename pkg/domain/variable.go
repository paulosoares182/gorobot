package domain

type VariableTemplate struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

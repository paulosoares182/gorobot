package domain

type Variable struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

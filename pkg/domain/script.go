package domain

type Script struct {
	ID       string    `json:"id"`
	Schedule string    `json:"schedule,omitempty"`
	Commands []Command `json:"commands"`
}

func NewScript() *Script {
	return &Script{
		ID:       "",
		Commands: []Command{},
	}
}

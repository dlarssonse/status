package status

// Status ...
type Status struct {
	Name     string
	Services []Service
}

// Service ...
type Service struct {
	Name     string
	Type     string
	Filename string
	Status   string
	Error    string `json:"error,omitempty"`
}

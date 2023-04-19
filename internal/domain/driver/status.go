package driver

type Status struct {
	stopped bool
}

func NewStatus(stopped bool) *Status {
	return &Status{stopped: stopped}
}

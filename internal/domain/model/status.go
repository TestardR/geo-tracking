package model

type Status struct {
	stopped bool
}

func NewStatus(stopped bool) Status {
	return Status{stopped: stopped}
}

func (s *Status) Stopped() bool {
	return s.stopped
}

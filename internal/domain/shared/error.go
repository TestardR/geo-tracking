package shared

type DomainError interface {
	Error() string
}

type DomainViolation struct {
	violation string
}

func (v DomainViolation) Error() string {
	return v.violation
}

func NewDomainError(violation string) DomainError {
	return DomainViolation{
		violation: violation,
	}
}

func IsDomainError(err error) bool {
	switch err.(type) {
	case DomainError:
		return true
	}
	return false
}

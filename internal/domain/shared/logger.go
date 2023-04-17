package shared

type ErrorInfoLogger interface {
	ErrorLogger
	InfoLogger
}

type DebugLogger interface {
	Debug(args ...interface{})
}

type ErrorLogger interface {
	Error(args ...interface{})
}

type InfoLogger interface {
	Info(args ...interface{})
}

type InfoFLogger interface {
	Infof(template string, args ...interface{})
}

type StdLogger interface {
	Printf(format string, v ...interface{})
}

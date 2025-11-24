package log

// public funcs

func Debug(msg string, args ...any) {
	logSingle().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logSingle().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logSingle().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logSingle().Error(msg, args...)
}

// writer
var LogWriter = writerSingle

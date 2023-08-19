package app

import "context"

type HTTPServer interface {
	ListenAndServe() error
	RegisterOnShutdown(f func())
	Shutdown(ctx context.Context) error
	Close() error
}

// Logger interface provides full minimal necessary list to log data
type Logger interface {
	// WithContext necessary to update Logger entity with any useful information from the context.Context
	WithContext(ctx context.Context) Logger
	// WithField adds any value into the Logger entity with the given name
	WithField(name string, value interface{}) Logger
	// WithError adds an error into the Logger entity with the given name
	WithError(err error) Logger

	// Debug log the message with the Debug level
	Debug(message string)
	// Info log the message with the Debug level
	Info(message string)
	// Warn log the message with the Warn level
	Warn(message string)
	// Error log the message with the Debug level
	Error(message string)
	// Fatal log the message with the Debug level
	Fatal(message string)
	// Panic log the message with the Debug level
	Panic(message string)
}

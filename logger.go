package app

import (
	"context"
	"os"
)

type noopLogger struct{}

var _ Logger = (*noopLogger)(nil)

func (l *noopLogger) WithContext(_ context.Context) Logger     { return l }
func (l *noopLogger) WithField(_ string, _ interface{}) Logger { return l }
func (l *noopLogger) WithError(_ error) Logger                 { return l }

func (l *noopLogger) Debug(_ string) {}
func (l *noopLogger) Info(_ string)  {}
func (l *noopLogger) Warn(_ string)  {}
func (l *noopLogger) Error(_ string) {}
func (l *noopLogger) Panic(msg string) {
	panic(msg)
}
func (l *noopLogger) Fatal(_ string) {
	os.Exit(1)
}

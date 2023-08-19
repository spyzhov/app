package app

import (
	"io"

	"go.uber.org/zap"
)

type closer struct {
	name   string
	closer io.Closer
}

// Close all necessary resources
func (app *Application) Close() {
	app.logger.Debug("Application stops")
	if app == nil {
		return
	}

	defer close(app.signal)
	defer close(app.error)

	for _, c := range app.closers {
		closer := c
		defer app.sClose(closer.closer, closer.name)
	}
}

// sClose will close any io.Closer object if it's not nil, any error will be only logged.
func (app *Application) sClose(closer io.Closer, scope string) {
	if !isNil(closer) {
		if err := closer.Close(); err != nil {
			app.logger.Warn("error on close", zap.String("scope", scope), zap.Error(err))
		}
	}
}

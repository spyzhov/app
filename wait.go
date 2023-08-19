package app

import (
	"fmt"

	"go.uber.org/zap"
)

// Wait waits while user don't press Ctrl+C or any error happens
func (app *Application) Wait() {
	defer app.Stop()

	select {
	case err := <-app.error:
		app.logger.Error("application error", zap.Error(err))
		panic(fmt.Errorf("application error: %w", err))
	case <-app.ctx.Done():
		app.logger.Info("stops via context", zap.Error(app.ctx.Err()))
	case sig := <-app.signal:
		app.logger.Info("stop", zap.Stringer("signal", sig))
	}
}

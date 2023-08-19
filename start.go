package app

import (
	"context"

	"go.uber.org/zap"
)

type runner struct {
	name     string
	function func() error
}

// Start initialize all long-living processes
func (app *Application) Start() *Application {
	app.logger.Debug("going to start application...")

	for _, r := range app.runners {
		if err := r.function(); err != nil {
			app.logger.Panic("runner start error", zap.String("name", r.name), zap.Error(err))
		}
	}

	return app
}

// Stop waits for all resources be cosed
func (app *Application) Stop() {
	if app == nil {
		return
	}
	app.logger.Info("application stops...")
	app.cancel()
	ctx, cancel := context.WithTimeout(context.Background(), app.appTimeout)

	go func() {
		defer cancel()
		app.wg.Wait()
	}()

	for _, fn := range app.onStop {
		fn(app.ctx)
	}

	<-ctx.Done()

	if ctx.Err() != context.Canceled {
		app.logger.Panic("service stopped with timeout")
	} else {
		app.logger.Info("service stopped with success")
	}
}

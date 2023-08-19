package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Application struct {
	// region System
	logger *zap.Logger
	info   BuildInfo
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	error chan error
	once  sync.Once

	signal chan os.Signal
	// endregion
	// region External
	healthchecks map[string]healthcheck
	hcTimeout    time.Duration // timeout for healthchecks

	runners []runner

	closers []closer

	onInit []func(context.Context) error
	onStop []func(context.Context)

	appTimeout time.Duration // timeout for Application to stop
	// endregion
}

func New(option ...Option) (app *Application, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app = &Application{
		// region System
		logger: zap.NewNop(),
		ctx:    ctx,
		cancel: cancel,
		error:  make(chan error, 1),
		wg:     new(sync.WaitGroup),
		signal: sigs,
		// endregion
		// region External
		healthchecks: make(map[string]healthcheck),
		hcTimeout:    time.Second,

		runners: make([]runner, 0),

		closers: make([]closer, 0),

		appTimeout: 3 * time.Second,

		onInit: make([]func(context.Context) error, 0),
		onStop: make([]func(context.Context), 0),
		// endregion
	}

	defer func() {
		if err != nil {
			app.Close()
		}
	}()

	for _, o := range option {
		o.apply(app)
	}

	return app, app.init()
}

// Error - register global error, but only once
func (app *Application) Error(err error) {
	if !isNil(err) {
		app.once.Do(func() {
			app.error <- err
		})
	}
}

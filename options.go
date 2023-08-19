package app

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (app *Application) With(options ...Option) *Application {
	Options(options).apply(app)
	return app
}

type Option interface {
	apply(app *Application)
}

type optionFunc func(app *Application)

func (f optionFunc) apply(app *Application) {
	f(app)
}

type Options []Option

func (fs Options) apply(app *Application) {
	for _, f := range fs {
		f.apply(app)
	}
}

// WithLogger set up a logger for Application
func WithLogger(log *zap.Logger) Option {
	return optionFunc(func(app *Application) {
		app.logger = log
	})
}

func WithInfoService(value string) Option {
	return optionFunc(func(app *Application) {
		app.info.Service = value
	})
}

func WithInfoVersion(value string) Option {
	return optionFunc(func(app *Application) {
		app.info.Version = value
	})
}

func WithInfoCreated(value string) Option {
	return optionFunc(func(app *Application) {
		app.info.Created = value
	})
}

func WithInfoCommit(value string) Option {
	return optionFunc(func(app *Application) {
		app.info.Commit = value
	})
}

func WithApplicationTimeout(value time.Duration) Option {
	return optionFunc(func(app *Application) {
		app.appTimeout = value
	})
}

func WithHealthcheck(name string, value func(ctx context.Context) error) Option {
	return optionFunc(func(app *Application) {
		app.healthchecks[name] = value
	})
}

func WithHealthcheckTimeout(value time.Duration) Option {
	return optionFunc(func(app *Application) {
		app.hcTimeout = value
	})
}

func WithCloser(name string, value io.Closer) Option {
	return optionFunc(func(app *Application) {
		app.closers = append(app.closers, closer{
			name:   name,
			closer: value,
		})
	})
}

func WithOnInit(value func(context.Context) error) Option {
	return optionFunc(func(app *Application) {
		app.onInit = append(app.onInit, value)
	})
}

func WithOnStop(value func(context.Context)) Option {
	return optionFunc(func(app *Application) {
		app.onStop = append(app.onStop, value)
	})
}

func WithHTTPHandler(name string, handler http.Handler, port int) Option {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: handler,
	}
	return WithHTTPServer(name, server, port)
}

func WithHTTPServer(name string, server HTTPServer, port int) Option {
	options := make(Options, 0)
	options = append(options, WithCloser("HTTPServer."+name, server))
	options = append(options, optionFunc(func(app *Application) {
		app.runners = append(app.runners, runner{
			name: name,
			function: func() error {
				log := app.logger.With(zap.String("name", name), zap.Int("port", port))

				app.wg.Add(1)
				go func() {
					defer app.wg.Done()
					defer func() {
						err := server.Shutdown(context.Background())
						if err != nil {
							log.Error("http server shutdown error", zap.Error(err))
						}
					}()
					server.RegisterOnShutdown(app.cancel)

					app.wg.Add(1)
					go func() {
						defer app.wg.Done()
						app.Error(server.ListenAndServe())
						log.Debug("http server not listen")
					}()

					log.Info("http server started")

					<-app.ctx.Done()

					log.Debug("http stops")
				}()
				return nil
			},
		})
	}))
	return options
}

func WithRunner(name string, fn func(ctx context.Context) error) Option {
	return optionFunc(func(app *Application) {
		app.runners = append(app.runners, runner{
			name: name,
			function: func() error {
				log := app.logger.With(zap.String("name", name))

				app.wg.Add(1)
				go func() {
					defer app.wg.Done()
					log.Info("runner starts")
					app.Error(fn(app.ctx))
					log.Debug("runner stops")
				}()
				return nil
			},
		})
	})
}

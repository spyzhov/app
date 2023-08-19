package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type healthcheck func(ctx context.Context) error

// Healthcheck is a handle function for health-check
func (app *Application) Healthcheck(ctx context.Context) (code int, info map[string]interface{}) {
	var mu sync.Mutex
	tests := make(map[string]string)
	code = http.StatusOK
	ctx, cancel := context.WithTimeout(ctx, app.hcTimeout)
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	for k, v := range app.healthchecks {
		key, validate := k, v

		g.Go(func() (err error) {
			defer func() {
				log := app.logger.With(zap.String("key", key))

				if r := recover(); r != nil {
					log = log.With(zap.Stack("stack"))
					err = fmt.Errorf("recovered: %v", r)
				}

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					tests[key] = "Error"
					log.Error("healthcheck:fail", zap.Error(err))
				} else {
					tests[key] = "OK"
					log.Debug("healthcheck:success")
				}
			}()

			return validate(ctx)
		})
	}

	if err := g.Wait(); err != nil {
		code = http.StatusInternalServerError
	}

	return code, map[string]interface{}{
		"info":  app.info,
		"time":  time.Now().Format(time.RFC3339),
		"tests": tests,
	}
}

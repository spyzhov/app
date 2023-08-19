package main

import (
	"context"
	"fmt"
	"math/rand"

	"go.uber.org/zap"

	"github.com/spyzhov/app"
)

func main() {
	// ... init ...
	application, err := app.New(
		app.WithInfoService("example"),
		app.WithInfoVersion("local"),

		app.WithLogger(zap.NewExample()),

		app.WithHealthcheck("random", func(ctx context.Context) error {
			if rand.Int()%2 == 0 {
				return fmt.Errorf("error")
			}
			return nil
		}),
		app.WithHealthcheck("success", func(ctx context.Context) error {
			return nil
		}),
	)
	application.
		With(
			app.WithHTTPHandler("management", app.NewManagementMux(application), 8080),
		)

	if err != nil {
		panic(err)
	}
	defer application.Close()
	application.
		Start().
		Wait()
}

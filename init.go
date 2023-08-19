package app

import (
	"fmt"
)

// init all necessary resources
func (app *Application) init() (err error) {
	for _, fn := range app.onInit {
		if err = fn(app.ctx); err != nil {
			return fmt.Errorf("cannot initialize application: %w", err)
		}
	}
	return nil
}

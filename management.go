package app

import (
	"encoding/json"
	"net/http"
)

func NewManagementMux(app *Application) *http.ServeMux {
	mux := http.NewServeMux()

	write := func(writer http.ResponseWriter, code int, info interface{}) {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(code)
		err := json.NewEncoder(writer).Encode(info)
		if err != nil {
			app.logger.Error("error on write response")
		}
	}

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		write(writer, http.StatusOK, app.Info())
	})
	mux.HandleFunc("/healthcheck", func(writer http.ResponseWriter, request *http.Request) {
		code, info := app.Healthcheck(request.Context())
		write(writer, code, info)
	})

	return mux
}

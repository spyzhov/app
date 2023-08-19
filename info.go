package app

type BuildInfo struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
	Created string `json:"created,omitempty"`
	Commit  string `json:"commit,omitempty"`
}

func (app *Application) Info() BuildInfo {
	return app.info
}

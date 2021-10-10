package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "environment": %q, "version": %q} `
	js = fmt.Sprintf(js, app.config.env, version)

	w.Header().Set("Content-Type", "applicaiton/json")

	w.Write([]byte(js))
}

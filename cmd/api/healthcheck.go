package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{ 
		"status":      "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil{
		app.logger.Println(err)
		msg := "The server encountered a problem and could not process your request"
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

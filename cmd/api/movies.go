package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Created Movie!")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "movie1, 2, 3")
}

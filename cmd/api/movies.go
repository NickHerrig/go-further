package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.nickherrig.com/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Created a Movie!")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Beetlejuice",
		Runtime:   92,
		Genres:    []string{"comedy", "fantasy"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		msg := "The server encountered a problem and could not process your request"
		http.Error(w, msg, http.StatusInternalServerError)
	}

}

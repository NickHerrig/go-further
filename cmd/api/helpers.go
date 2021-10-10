package main

import (
	"net/http"
	"strconv"
	"errors"

	"github.com/julienschmidt/httprouter"
)


func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64) //base 10, 64 bit
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

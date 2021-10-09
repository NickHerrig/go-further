package main

import (
	"fmt"
	"log"
)

const version = "0.0.1"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	loger *log.Logger
}

func main() {
	var cfg config
	fmt.Println(cfg)
	fmt.Println("hello world!")
	// implement cli params

	// implement logger

	// implement app struct

	// create a server mux
	// TODO Read about servermux godoc
	// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

	// TODO Read up on http.Server
	//https://pkg.go.dev/net/http


	// start the go server
}

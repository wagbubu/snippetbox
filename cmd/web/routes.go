package main

import (
	"net/http"

	"github.com/wagbubu/snippetbox/middlewares"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(middlewares.NeuteredFileSystem{Fs: http.Dir("./ui/static/")})

	mux.HandleFunc("/", app.home)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}

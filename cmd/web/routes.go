package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{Fs: http.Dir("./ui/static/")})

	mux.HandleFunc("/", app.home)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.recoverPanic((app.logRequest(secureHeaders(mux))))
}

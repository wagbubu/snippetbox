package main

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	Fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	/*
		Opens file, we will return the file if the path
		is not a directory like /static/main.css or /static/index.html
	*/
	f, err := nfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}
	/*
		Get file stats/metadata/information
		to know if its a directory or not
	*/
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	/*
		check if the file is a directory, if it is
		check if it has index.html
		if not we will return error
		which will be a not found error
	*/
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.Fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	/*
		if the directory has index then we will utilize Go's autmatic serving
		of index.html if it exists in a directory
		the above condition is just to avoid showing all files in a directory
		that doesnt have index.html
	*/
	return f, nil
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' cdn.jsdelivr.net; style-src 'self' fonts.googleapis.com cdn.jsdelivr.net; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

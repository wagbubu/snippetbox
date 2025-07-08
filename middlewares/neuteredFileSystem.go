package middlewares

import (
	"net/http"
	"path/filepath"
)

type NeuteredFileSystem struct {
	Fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
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

package web

import (
	"net/http"
	"os"
	"path/filepath"
)

func Server(w http.ResponseWriter, r *http.Request) {

	path := filepath.Join("build", r.URL.Path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join("build", "index.html"))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir("build")).ServeHTTP(w, r)
}

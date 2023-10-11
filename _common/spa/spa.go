package spa

import (
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct {
	Static string
	Index  string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.Static, r.URL.Path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.Static, h.Index))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Cache-Control", "no-cache")
	http.FileServer(http.Dir(h.Static)).ServeHTTP(w, r)
}

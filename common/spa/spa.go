package spa

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
	Static    string
	Index     string
	ServeGzip bool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.Static, r.URL.Path)
	fileType := filepath.Ext(path)

	acceptGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	// Check whether a file exists or is a directory
	if fi, err := os.Stat(path); os.IsNotExist(err) || fi.IsDir() {
		htmlFile := filepath.Join(h.Static, h.Index)
		if h.ServeGzip && acceptGzip {
			if _, err := os.Stat(htmlFile + ".gz"); err == nil {
				w.Header().Add("Content-Encoding", "gzip")
				w.Header().Add("Content-Type", "text/html")
				htmlFile += ".gz"
			}
		}

		// Serve SPA
		http.ServeFile(w, r, htmlFile)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve gziped
	if h.ServeGzip && acceptGzip && (fileType == ".js" || fileType == ".css") {
		if _, err := os.Stat(path + ".gz"); err == nil {

			w.Header().Add("Content-Encoding", "gzip")
			if fileType == ".js" {
				w.Header().Add("Content-Type", "text/javascript")
			}
			if fileType == ".css" {
				w.Header().Add("Content-Type", "text/css")
			}

			http.ServeFile(w, r, filepath.Join(path+".gz"))
			return
		}
	}

	w.Header().Add("Cache-Control", "no-cache")
	http.FileServer(http.Dir(h.Static)).ServeHTTP(w, r)
}

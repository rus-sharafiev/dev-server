package fswr

import (
	"net/http"
	"os"
	"path"
	"strings"
)

func FileServerWithRedirect(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		// attempt to open the file via the http.FileSystem
		file, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				r.URL.Path = "/"
			}
		}

		if err == nil {
			file.Close()
		}

		http.FileServer(root).ServeHTTP(w, r)
	})
}

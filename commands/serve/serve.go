package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rus-sharafiev/dev/common"
	"github.com/rus-sharafiev/dev/common/browser"
	"github.com/rus-sharafiev/dev/common/spa"
)

func Run() {

	// Web server
	router := http.NewServeMux()
	router.Handle("/", spa.Handler{
		Static:    "build",
		Index:     "index.html",
		ServeGzip: true,
	})

	port := "8000"
	if common.Config.Port != nil {
		port = *common.Config.Port
	}

	fmt.Printf("\n\x1b[2mHTTP server is running on http://localhost:%v/\n \x1b[0m ", port)

	go browser.Open("http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

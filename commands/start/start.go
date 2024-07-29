package start

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
	"github.com/rus-sharafiev/dev/_common/browser"
	"github.com/rus-sharafiev/dev/_common/spa"
	"github.com/rus-sharafiev/dev/plugins/less"
	"github.com/rus-sharafiev/dev/plugins/sass"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients []*websocket.Conn

var reloadPlugin = api.Plugin{
	Name: "reload-plugin",
	Setup: func(build api.PluginBuild) {

		build.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
			for _, conn := range clients {
				conn.WriteMessage(websocket.TextMessage, []byte("reload"))
			}
			clients = nil
			return api.OnEndResult{}, nil
		})

	},
}

func Run() {

	// esbuild
	ctx, err := api.Context(api.BuildOptions{
		EntryPoints: []string{"src/*.ts*", "src/index.html"},
		JSXDev:      true,
		JSX:         api.JSXAutomatic,
		Bundle:      true,
		Outdir:      "build",
		Sourcemap:   api.SourceMapLinked,
		Plugins:     []api.Plugin{reloadPlugin, less.Plugin, sass.Plugin},
		External: []string{
			"*.gif",
			"*.svg",
			"*.jpg",
			"*.png",
		},
		Loader: map[string]api.Loader{
			".woff":  api.LoaderDataURL,
			".woff2": api.LoaderDataURL,
			".otf":   api.LoaderDataURL,
			".eot":   api.LoaderDataURL,
			".ttf":   api.LoaderDataURL,
			".html":  api.LoaderCopy,
		},
		Banner:   map[string]string{"js": "(() => new WebSocket('ws://localhost:8000/ws').onmessage = () => location.reload())(); var isDevBuild = true;"},
		Write:    true,
		LogLevel: api.LogLevelInfo,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := ctx.Watch(api.WatchOptions{}); err != nil {
		log.Fatal(err)
	}

	// Web server
	router := http.NewServeMux()
	router.Handle("/", spa.Handler{
		Static: "build",
		Index:  "index.html",
	})

	// Live reload via websocket
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		clients = append(clients, conn)
	})

	fmt.Printf("\n\x1b[2mHTTP server is running on http://localhost:8000/\n \x1b[0m ")
	fmt.Printf("\n\x1b[33m[esbuild] \x1b[0mwatching for changes...\n\n")

	go browser.Open("http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", router))
}

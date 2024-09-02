package start

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
	"github.com/rus-sharafiev/dev/_common/browser"
	"github.com/rus-sharafiev/dev/_common/conf"
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

func Run(conf *conf.DevConfig) {

	entryPoints := []string{"src/index.ts*"}
	cssLoader := api.LoaderCSS
	port := "8000"
	bundle := true
	external := []string{
		"*.gif",
		"*.svg",
		"*.jpg",
		"*.png",
	}

	if conf != nil {
		if conf.EntryPoints != nil {
			entryPoints = *conf.EntryPoints
		}
		if conf.Port != nil {
			port = *conf.Port
		}
		if conf.Bundle != nil {
			bundle = *conf.Bundle

			if !bundle {
				external = nil
			}
		}
		if conf.WebComponents != nil {
			cssLoader = api.LoaderText
		}
	}

	// esbuild
	ctx, err := api.Context(api.BuildOptions{
		EntryPoints: entryPoints,
		JSXDev:      true,
		JSX:         api.JSXAutomatic,
		Bundle:      bundle,
		Outdir:      "build",
		Sourcemap:   api.SourceMapLinked,
		Plugins:     []api.Plugin{reloadPlugin, less.Plugin, sass.Plugin},
		External:    external,
		Loader: map[string]api.Loader{
			".woff":  api.LoaderDataURL,
			".woff2": api.LoaderDataURL,
			".otf":   api.LoaderDataURL,
			".eot":   api.LoaderDataURL,
			".ttf":   api.LoaderDataURL,
			".html":  api.LoaderCopy,
			".css":   cssLoader,
		},
		Banner:   map[string]string{"js": "(() => new WebSocket('ws://localhost:" + port + "/ws').onmessage = () => location.reload())(); var isDevBuild = true;"},
		Write:    true,
		LogLevel: api.LogLevelInfo,
		Engines: []api.Engine{
			{Name: api.EngineChrome, Version: "100"},
		},
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
		Static:    "build",
		Index:     "index.html",
		ServeGzip: false,
	})

	// Live reload via websocket
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		clients = append(clients, conn)
	})

	fmt.Printf("\n\x1b[2mHTTP server is running on http://localhost:%v/\n \x1b[0m ", port)
	fmt.Printf("\n\x1b[33m[esbuild] \x1b[0mwatching for changes...\n\n")

	go browser.Open("http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

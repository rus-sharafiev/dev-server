package start

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
	"github.com/rus-sharafiev/dev/common"
	"github.com/rus-sharafiev/dev/common/browser"
	"github.com/rus-sharafiev/dev/common/spa"
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

	entryPoints := []string{"src/index.ts*"}
	cssLoader := api.LoaderCSS
	port := "8000"
	bundle := true
	charset := api.CharsetUTF8
	external := []string{
		"*.gif",
		"*.svg",
		"*.jpg",
		"*.png",
	}

	if common.Config.External != nil {
		configExternal := *common.Config.External
		external = append(external, configExternal...)
	}

	if common.Config.EntryPoints != nil {
		entryPoints = *common.Config.EntryPoints
	}
	if common.Config.Port != nil {
		port = *common.Config.Port
	}
	if common.Config.Bundle != nil {
		bundle = *common.Config.Bundle

		if !bundle {
			external = nil
		}
	}

	switch common.Config.Charset {
	case "default":
		charset = api.CharsetDefault
	case "ascii":
		charset = api.CharsetASCII
	}

	format := api.FormatESModule
	switch common.Config.Format {
	case "iife":
		format = api.FormatIIFE
	case "cjs":
		format = api.FormatCommonJS
	case "default":
		format = api.FormatDefault
	}

	target := api.ES2020
	switch common.Config.Target {
	case "ES2018":
		target = api.ES2018
	case "ES2022":
		target = api.ES2022
	}

	if common.Config.WebComponents != nil {
		cssLoader = api.LoaderText
		format = api.FormatDefault
	}

	// esbuild
	ctx, err := api.Context(api.BuildOptions{
		EntryPoints: entryPoints,
		JSXDev:      true,
		JSX:         api.JSXAutomatic,
		Bundle:      bundle,
		Outdir:      "build",
		Charset:     charset,
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
		Format:   format,
		Target:   target,
		Engines: []api.Engine{
			{Name: api.EngineChrome, Version: "100"},
			{Name: api.EngineFirefox, Version: "100"},
			{Name: api.EngineSafari, Version: "15"},
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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rus-sharafiev/dev-server/fswr"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func main() {
	ctx, err := api.Context(api.BuildOptions{
		EntryPoints: []string{"src/index.tsx"},
		JSXDev:      true,
		JSX:         api.JSXAutomatic,
		Bundle:      true,
		Outdir:      "build",
		Sourcemap:   api.SourceMapInline,
		Plugins:     []api.Plugin{reloadPlugin},
		External:    []string{"*.gif", "*.eot", "*.woff", "*.ttf"},
		Banner:      map[string]string{"js": "(() => new WebSocket('ws://localhost:8000/ws').onmessage = () => location.reload())();"},
		Write:       true,
		LogLevel:    api.LogLevelInfo,
	})

	if err != nil {
		os.Exit(1)
	}

	err2 := ctx.Watch(api.WatchOptions{})
	if err2 != nil {
		os.Exit(1)
	}

	fmt.Printf("\n\x1b[2mHTTP server is running on http://localhost:8000/\n \x1b[0m ")
	fmt.Printf("\n\x1b[33m[esbuild] \x1b[0mwatching for changes...\n\n")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		clients = append(clients, conn)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.Handle("/", http.StripPrefix("/", fswr.FileServerWithRedirect(http.Dir("build/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

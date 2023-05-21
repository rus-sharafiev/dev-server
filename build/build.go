package build

import (
	"os"
	lessplugin "rus-sharafiev/dev-server/less-plugin"

	"github.com/evanw/esbuild/pkg/api"
)

func Run() {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{"src/index.tsx"},
		JSXDev:      true,
		JSX:         api.JSXAutomatic,
		Bundle:      true,
		Outdir:      "dist",
		Sourcemap:   api.SourceMapInline,
		Plugins:     []api.Plugin{lessplugin.LessPlugin},
		External:    []string{"*.gif", "*.eot", "*.woff", "*.ttf"},
		Write:       true,
		LogLevel:    api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}

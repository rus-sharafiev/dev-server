package build

import (
	"os"
	"rus-sharafiev/dev-server/less"

	"github.com/evanw/esbuild/pkg/api"
)

func Run() {
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"src/index.tsx"},
		JSXDev:            true,
		JSX:               api.JSXAutomatic,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Outdir:            "build",
		Sourcemap:         api.SourceMapInline,
		Plugins:           []api.Plugin{less.Plugin},
		External:          []string{"*.gif", "*.eot", "*.woff", "*.ttf"},
		Write:             true,
		LogLevel:          api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}

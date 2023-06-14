package build

import (
	"os"
	"rus-sharafiev/dev-server/less"
	"rus-sharafiev/dev-server/sass"

	"github.com/evanw/esbuild/pkg/api"
)

func Run() {
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"src/index.tsx"},
		JSX:               api.JSXAutomatic,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Outdir:            "build",
		Sourcemap:         api.SourceMapLinked,
		Plugins:           []api.Plugin{less.Plugin, sass.Plugin},
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
		},
		Write:    true,
		LogLevel: api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}

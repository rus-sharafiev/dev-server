package build

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rus-sharafiev/dev/plugins/less"
	"github.com/rus-sharafiev/dev/plugins/sass"
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

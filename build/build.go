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
		JSXDev:            true,
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
		},
		Loader: map[string]api.Loader{
			".jpg":   api.LoaderDataURL,
			".png":   api.LoaderDataURL,
			".woff":  api.LoaderDataURL,
			".woff2": api.LoaderDataURL,
			".otf":   api.LoaderDataURL,
			".eot":   api.LoaderDataURL,
			".ttf":   api.LoaderDataURL,
			".svg":   api.LoaderText,
		},
		Write:    true,
		LogLevel: api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}

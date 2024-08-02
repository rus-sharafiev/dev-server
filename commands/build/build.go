package build

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rus-sharafiev/dev/_common/conf"
	"github.com/rus-sharafiev/dev/plugins/less"
	"github.com/rus-sharafiev/dev/plugins/sass"
)

func Run(conf *conf.DevConfig) {

	entryPoints := []string{"src/*.ts*"}
	if conf != nil && conf.EntryPoints != nil {
		entryPoints = append(entryPoints, *conf.EntryPoints...)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:       entryPoints,
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
			".html":  api.LoaderCopy,
		},
		Write:    true,
		LogLevel: api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	for _, file := range result.OutputFiles {
		if fileType := filepath.Ext(file.Path); fileType == ".js" || fileType == ".css" || fileType == ".html" {

			var b bytes.Buffer
			gw := gzip.NewWriter(&b)
			gw.Write(file.Contents)
			gw.Close()

			if err := os.WriteFile(file.Path+".gz", b.Bytes(), 0666); err != nil {
				fmt.Println(err)
			}
		}
	}
}

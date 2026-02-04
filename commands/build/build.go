package build

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rus-sharafiev/dev/common"
	"github.com/rus-sharafiev/dev/plugins/less"
	"github.com/rus-sharafiev/dev/plugins/sass"
)

func Run() {

	entryPoints := []string{"src/index.ts*"}
	cssLoader := api.LoaderCSS
	bundle := true
	keepNames := false
	charset := api.CharsetUTF8
	external := []string{
		"*.gif",
		"*.svg",
		"*.jpg",
		"*.png",
	}

	var minifyCssErrors []api.Message = nil

	if common.Config.External != nil {
		configExternal := *common.Config.External
		external = append(external, configExternal...)
	}
	if common.Config.EntryPoints != nil {
		entryPoints = *common.Config.EntryPoints
	}
	if common.Config.Bundle != nil {
		bundle = *common.Config.Bundle

		if !bundle {
			keepNames = true
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
		minifyCssErrors = minifyCss()
	}

	if minifyCssErrors != nil {
		for _, err := range minifyCssErrors {
			fmt.Printf("\nWeb components: \x1b[31m%v: %v\x1b[0m", "Error preparing css files", err)
		}
		return
	}

	wcCssPlugin := api.Plugin{
		Name: "wc-css-plugin",
		Setup: func(build api.PluginBuild) {

			build.OnLoad(api.OnLoadOptions{Filter: `\.ts$`, Namespace: `file`},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					resolveDir := filepath.Dir(args.Path)

					b, err := os.ReadFile(args.Path)
					if err != nil {
						log.Fatal(err)
					}

					contents := strings.ReplaceAll(string(b), "/styles/", "/styles/_minified/")

					return api.OnLoadResult{
						Contents:   &contents,
						Loader:     api.LoaderTS,
						ResolveDir: resolveDir,
					}, nil
				})

		},
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:       entryPoints,
		JSX:               api.JSXAutomatic,
		Bundle:            bundle,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		KeepNames:         keepNames,
		Outdir:            "build",
		Charset:           charset,
		Sourcemap:         api.SourceMapLinked,
		Plugins:           []api.Plugin{less.Plugin, sass.Plugin, wcCssPlugin},
		External:          external,
		Loader: map[string]api.Loader{
			".woff":  api.LoaderDataURL,
			".woff2": api.LoaderDataURL,
			".otf":   api.LoaderDataURL,
			".eot":   api.LoaderDataURL,
			".ttf":   api.LoaderDataURL,
			".html":  api.LoaderCopy,
			".css":   cssLoader,
		},
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

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	if common.Config.CreateGzip != nil && *common.Config.CreateGzip {
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
}

func minifyCss() []api.Message {
	if common.Config.WebComponents.StylesDir != nil {
		result := api.Build(api.BuildOptions{
			EntryPoints:       []string{*common.Config.WebComponents.StylesDir + "/*.css"},
			Outdir:            *common.Config.WebComponents.StylesDir + "/_minified",
			Bundle:            true,
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
			Write:             true,
			Engines: []api.Engine{
				{Name: api.EngineChrome, Version: "100"},
			},
		})

		if len(result.Errors) > 0 {
			return result.Errors
		}

		return nil
	}
	return nil
}

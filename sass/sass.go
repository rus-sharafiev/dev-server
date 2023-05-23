package sass

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"
)

var Plugin = api.Plugin{
	Name: "sass-plugin",
	Setup: func(build api.PluginBuild) {

		build.OnResolve(api.OnResolveOptions{Filter: `.\.(scss|sass)$`, Namespace: `file`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {

				pathResolve := build.Resolve(args.Path, api.ResolveOptions{
					Kind:       args.Kind,
					Importer:   args.Importer,
					ResolveDir: args.ResolveDir,
					PluginData: args.PluginData,
				})

				filePath := pathResolve.Path

				watchFiles := []string{filePath}

				watchFiles = append(watchFiles, getSassImports(filePath)...)
				return api.OnResolveResult{
					Path:       filePath,
					WatchFiles: watchFiles,
				}, nil
			})

		build.OnLoad(api.OnLoadOptions{Filter: `.\.(scss|sass)$`, Namespace: `file`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				resolveDir := filepath.Dir(args.Path)

				cmd := exec.Command("sass", args.Path)

				out, err := cmd.CombinedOutput()
				if err != nil {
					return api.OnLoadResult{}, errors.New(string(out))
				}

				contents := string(out)

				return api.OnLoadResult{
					Contents:   &contents,
					Loader:     api.LoaderCSS,
					ResolveDir: resolveDir,
				}, nil
			})

	},
}

func getSassImports(filePath string) []string {
	rootDir := filepath.Dir(filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return []string{"Less plugin: Error reading @import file"}
	}

	importRegex := regexp.MustCompile(`@import.*?["']([^"']+)["'].*?`)
	match := importRegex.FindAllStringSubmatch(string(content), -1)

	if len(match) == 0 {
		return []string{}
	}

	result := []string{}

	for i := range match {
		dirName, fileName := filepath.Split(match[i][1])
		fullPath := filepath.Join(rootDir, filepath.Clean(dirName), fileName)
		result = append(result, fullPath)
		result = append(result, getSassImports(fullPath)...)
	}

	return result
}

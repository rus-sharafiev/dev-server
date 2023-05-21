package lessplugin

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"
)

var LessPlugin = api.Plugin{
	Name: "less-plugin",
	Setup: func(build api.PluginBuild) {

		// Resolve *.less files with namespace
		build.OnResolve(api.OnResolveOptions{Filter: `\.less$`, Namespace: `file`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {

				pathResolve := build.Resolve(args.Path, api.ResolveOptions{
					Kind:       args.Kind,
					Importer:   args.Importer,
					ResolveDir: args.ResolveDir,
					PluginData: args.PluginData,
				})

				filePath := pathResolve.Path

				watchFiles := []string{filePath}

				watchFiles = append(watchFiles, getLessImports(filePath)...)
				return api.OnResolveResult{
					Path:       filePath,
					WatchFiles: watchFiles,
				}, nil
			})

		// Build .less files
		build.OnLoad(api.OnLoadOptions{Filter: `\.less$`, Namespace: `file`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				resolveDir := filepath.Dir(args.Path)

				cmd := exec.Command(".\\less-compiler.exe", args.Path)

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

// Recursively get .less/.css imports from file
func getLessImports(filePath string) []string {
	rootDir := filepath.Dir(filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return []string{"Less plugin: Error reading @import file"}
	}

	importCommentRegex := regexp.MustCompile(`(?m)(?:\/\*(?:[\s\S]*?)\*\/)|(\/\/(?:.*)$)`)
	cleanContent := importCommentRegex.ReplaceAllString(string(content), string(""))

	importRegex := regexp.MustCompile(`@import.*?["']([^"']+)["'].*?`)
	match := importRegex.FindAllStringSubmatch(cleanContent, -1)

	if len(match) == 0 {
		return []string{}
	}

	result := []string{}

	for i := range match {
		dirName, fileName := filepath.Split(match[i][1])
		fullPath := filepath.Join(rootDir, filepath.Clean(dirName), fileName)
		result = append(result, fullPath)
		result = append(result, getLessImports(fullPath)...)
	}

	return result
}

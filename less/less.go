package less

import (
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tystuyfzand/less-go"
)

var LessPlugin = api.Plugin{
	Name: "less",
	Setup: func(build api.PluginBuild) {

		// Resolve *.less files with namespace
		build.OnResolve(api.OnResolveOptions{Filter: `/\.less$/`, Namespace: `file`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {

				pathResolve := build.Resolve(args.Path, api.ResolveOptions{
					Kind:       args.Kind,
					Importer:   args.Importer,
					ResolveDir: args.ResolveDir,
					PluginData: args.PluginData,
				})

				filePath := pathResolve.Path

				watchFiles := []string{filePath}

				return api.OnResolveResult{
					Path:       filePath,
					WatchFiles: watchFiles,
				}, nil
			})

		// Build .less files
		build.OnLoad(api.OnLoadOptions{Filter: `\.txt$`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				dir := filepath.Dir(args.Path)

				content, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				opts := map[string]interface{}{
					"filename":     args.Path,
					"compress":     true,
					"relativeUrls": true,
					"paths":        []string{args.Path},
				}

				contents, err := less.Render(string(content), opts)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				return api.OnLoadResult{
					Contents:   &contents,
					Loader:     api.LoaderCSS,
					ResolveDir: dir,
				}, nil
			})

	},
}

// const importRegex = `/@import.*?["']([^"']+)["'].*?/`
// const globalImportRegex = `/@import.*?["']([^"']+)["'].*?/g`
// const importCommentRegex = `/(?:\/\*(?:[\s\S]*?)\*\/)|(\/\/(?:.*)$)/gm`

// var extWhitelist = [2]string{".css", ".less"}

// Recursively get .less/.css imports from file
// func getLessImports(filePath string, paths []string) []string {
// 	dir := filepath.Dir(filePath)

// 	content, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return []string{}
// 	}

// 	cleanContent := strings.Replace()
// }

// Convert less error into esbuild error
// func convertLessError(err less.Compiler.error)

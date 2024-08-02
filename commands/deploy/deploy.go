package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/rus-sharafiev/dev/_common/conf"
)

func Run(conf *conf.DevConfig) {

	if conf == nil {
		fmt.Printf("\nConfig file: \x1b[31m%v\x1b[0m\n\n", "Deploy config hasn't been provided")
		return
	}

	entries, err := os.ReadDir("build")
	if err != nil {
		fmt.Printf("\nError reading project build directory: \x1b[31m%v\x1b[0m\n\n", err)
		return
	}

	if conf.CssPath != nil && conf.JsPath != nil {
		fmt.Printf("\n%s\n\n", "Copying files via scp...")

		var jsFiles []string
		var cssFiles []string

		jsRe := regexp.MustCompile(`^.*\.(js|js\.gz|js\.map)$`)
		cssRe := regexp.MustCompile(`^.*\.(css|css\.gz|css\.map)$`)

		for _, e := range entries {
			if jsRe.MatchString(e.Name()) {
				jsFiles = append(jsFiles, "build/"+e.Name())
			} else if cssRe.MatchString(e.Name()) {
				cssFiles = append(cssFiles, "build/"+e.Name())
			}
		}

		copyViaScp(jsFiles, *conf.JsPath)
		copyViaScp(cssFiles, *conf.CssPath)

	} else if confPath := conf.DeployPath; confPath != nil {
		fmt.Printf("\n%s\n\n", "Copying files via scp...")

		var files []string
		for _, e := range entries {
			files = append(files, "build/"+e.Name())
		}

		copyViaScp(files, *conf.DeployPath)

	} else {
		fmt.Printf("\n\x1b[31mConfig file:\x1b[0m %v\n\n", "Deploy path hasn't been provided")
		return
	}

	fmt.Printf("\n\x1b[32m%s\x1b[0m\n\n", "Successfully copied")
}

func copyViaScp(files []string, path string) {
	if len(files) == 0 {
		return
	}

	commandArgs := []string{"-r"}
	commandArgs = append(commandArgs, files...)
	commandArgs = append(commandArgs, path)

	cmd := exec.Command("scp", commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

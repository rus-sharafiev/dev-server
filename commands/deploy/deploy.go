package deploy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const confFile = "./dev.conf"

type DevConfig struct {
	DeployPath *string `json:"deployPath"`
	JsPath     *string `json:"jsPath"`
	CssPath    *string `json:"cssPath"`
}

func Run(args ...string) {

	entries, err := os.ReadDir("build")
	if err != nil {
		fmt.Printf("\n\x1b[31mReadDir err:\x1b[0m %v\n", err)
	}

	fmt.Printf("\n%s\n\n", "Copying files via scp...")

	if len(args) >= 1 {

		var files []string
		for _, e := range entries {
			files = append(files, "build/"+e.Name())
		}

		copyViaScp(files, args[0])

	} else {

		if _, err := os.Stat(confFile); err != nil && os.IsNotExist(err) {
			log.Fatal("No target path has been provided, e.g. dev deploy root@0.0.0.0:/var/www/html/")
		} else if err != nil {
			log.Fatalf("\n\x1b[31mConfig file:\x1b[0m %v\n", "Error reading the file")
		}

		data, err := os.ReadFile(confFile)
		if err != nil {
			log.Fatalf("\n\x1b[31mReadFile err:\x1b[0m %v\n", err)
		}

		var config DevConfig
		if err = json.Unmarshal(data, &config); err != nil {
			log.Fatalf("\n\x1b[31mUnmarshal err:\x1b[0m %v\n", err)
		}

		if config.CssPath != nil && config.JsPath != nil {

			var jsFiles []string
			var cssFiles []string

			for _, e := range entries {
				if fileType := filepath.Ext(e.Name()); fileType == ".js" {
					jsFiles = append(jsFiles, "build/"+e.Name())
				} else if fileType == ".css" {
					cssFiles = append(cssFiles, "build/"+e.Name())
				}
			}

			copyViaScp(jsFiles, *config.JsPath)
			copyViaScp(cssFiles, *config.CssPath)

		} else if confPath := config.DeployPath; confPath != nil {

			var files []string
			for _, e := range entries {
				files = append(files, "build/"+e.Name())
			}

			copyViaScp(files, *config.DeployPath)

		} else {
			log.Fatalf("\n\x1b[31mConfig file:\x1b[0m %v\n", "no deploy path has been provided")
		}
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

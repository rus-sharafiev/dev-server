package deploy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const confFile = "./dev.conf"

type DevConfig struct {
	DeployPath string `json:"deployPath"`
}

func Run(args ...string) {
	var path string

	if len(args) >= 1 {
		path = args[0]
	} else {
		if _, err := os.Stat(confFile); os.IsNotExist(err) {

			log.Fatal("No target path has been provided, e.g. dev deploy root@0.0.0.0:/var/www/html/")

		} else if err == nil {

			data, err := os.ReadFile(confFile)
			if err != nil {
				log.Fatal("ReadFile err: ", err)
			}

			var config DevConfig
			if err = json.Unmarshal(data, &config); err != nil {
				log.Fatal("Unmarshal err: ", err)
			}

			path = config.DeployPath
		}
	}

	entries, err := os.ReadDir("build")
	if err != nil {
		log.Fatal(err)
	}

	commandArgs := []string{"-r"}
	for _, e := range entries {
		commandArgs = append(commandArgs, "build/"+e.Name())
	}
	commandArgs = append(commandArgs, path)

	fmt.Printf("\n%s\n\n", "Copying files via scp...")

	cmd := exec.Command("scp", commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\x1b[32m%s\x1b[0m\n\n", "Successfully copied")
}

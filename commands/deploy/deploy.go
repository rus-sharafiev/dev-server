package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const envFile = "dev.env"

func Run(args ...string) {
	var path string

	if len(args) >= 1 {
		path = args[0]
	} else {
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			log.Fatal("No target path has been provided, e.g. dev deploy root@0.0.0.0:/var/www/html/")
		} else if err == nil {
			data, err := os.ReadFile(envFile)
			if err != nil {
				log.Fatal(err)
			}
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				part := strings.Split(line, "=")

				if part[0] == "deployPath" {
					path = part[1]
				}
			}
		}
	}

	fmt.Println(path)

	fmt.Printf("\n%s\n\n", "Copying files via scp...")

	cmd := exec.Command("scp", "-r", "build/*", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\x1b[32m%s\x1b[0m\n\n", "Successfully copied")
}

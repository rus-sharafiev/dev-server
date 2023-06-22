package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(path string) {

	currentDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	buildDir := filepath.Join(currentDir, "build", "*")
	unixPath := strings.ReplaceAll(buildDir, "\\", "/")

	fmt.Printf("\n%s\n\n", "Copying files via scp...")

	cmd := exec.Command("scp", "-r", unixPath, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n\x1b[32m%s\x1b[0m\n\n", "Successfully copied")
}

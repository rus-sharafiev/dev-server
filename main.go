package main

import (
	"fmt"
	"os"
	"rus-sharafiev/dev-server/build"
	"rus-sharafiev/dev-server/start"
)

func main() {

	if noArgs := len(os.Args); noArgs == 1 {
		fmt.Println("No argument has been provided")
		return
	}

	switch script := os.Args[1]; script {
	case "start":
		start.Run()
	case "build":
		build.Run()
	default:
		fmt.Println("Invalid argument. Allowed are 'build' and 'dev'")
	}
}

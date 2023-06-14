package main

import (
	"fmt"
	"os"
	"rus-sharafiev/dev-server/build"
	"rus-sharafiev/dev-server/deploy"
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

	case "deploy":
		if argLength := len(os.Args); argLength == 3 {
			build.Run()
			deploy.Run(os.Args[2])
		} else {
			fmt.Println("No target path has been provided, e.g. root@0.0.0.0:/var/www/html/")
		}

	default:
		fmt.Println("Invalid argument. Allowed are 'build' and 'dev'")
	}
}

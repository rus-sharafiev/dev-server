package main

import (
	"fmt"
	"os"

	"github.com/rus-sharafiev/dev/commands/build"
	"github.com/rus-sharafiev/dev/commands/create"
	"github.com/rus-sharafiev/dev/commands/deploy"
	"github.com/rus-sharafiev/dev/commands/serve"
	"github.com/rus-sharafiev/dev/commands/start"
	"github.com/rus-sharafiev/dev/common"
)

func main() {

	// Load app config
	common.LoadConf()

	if noArgs := len(os.Args); noArgs == 1 {
		fmt.Println("No argument has been provided")
		return
	}

	switch script := os.Args[1]; script {

	case "start":
		start.Run()

	case "build":
		build.Run()

	case "serve":
		build.Run()
		serve.Run()

	case "deploy":
		build.Run()
		deploy.Run()

	case "create":
		if argLength := len(os.Args); argLength == 3 {

			switch arg := os.Args[2]; arg {

			case "project":
				create.Build()
				create.Src()

			case "build":
				create.Build()

			case "src":
				create.Src()

			default:
				fmt.Println("Invalid argument. Allowed are 'build' or 'src'")
			}

		} else {
			create.Build()
			create.Src()
		}

	default:
		fmt.Println("Invalid argument. Allowed are 'start', 'build', 'deploy' or 'create'")
	}
}

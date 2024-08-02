package main

import (
	"fmt"
	"os"

	"github.com/rus-sharafiev/dev/_common/conf"
	"github.com/rus-sharafiev/dev/commands/build"
	"github.com/rus-sharafiev/dev/commands/create"
	"github.com/rus-sharafiev/dev/commands/deploy"
	"github.com/rus-sharafiev/dev/commands/serve"
	"github.com/rus-sharafiev/dev/commands/start"
)

func main() {

	if noArgs := len(os.Args); noArgs == 1 {
		fmt.Println("No argument has been provided")
		return
	}

	conf := conf.Get()

	switch script := os.Args[1]; script {

	case "start":
		start.Run(conf)

	case "build":
		build.Run(conf)

	case "serve":
		build.Run(conf)
		serve.Run(conf)

	case "deploy":
		build.Run(conf)
		deploy.Run(conf)

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

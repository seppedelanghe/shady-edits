package main

import (
	"os"
	"shady-edits/pkg/app"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		panic("need at least 2 file paths to load")
	}

	config, err := app.NewConfigFromPaths(args[0], args[1])
	if err != nil {
		panic(err)
	}

	debugApp := app.NewDebugApp(config)
	if err = debugApp.Run(); err != nil {
		panic(err)
	}

}

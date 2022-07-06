package main

import (
	"os"
	"runtime"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	var goblEnv []string
	var extraArgs []interface{}
	if len(os.Args) > 2 {
		split := len(os.Args) - 3
		for i, a := range os.Args[2:] {
			if a == "--" {
				split = i
				break
			}
		}
		goblEnv = os.Args[2 : 2+split]
		for _, v := range os.Args[2+split+1:] {
			extraArgs = append(extraArgs, v)
		}
	}

	// Adjust exe if an env is GOOD=windows
	for _, v := range goblEnv {
		if v == "GOOS=windows" {
			exe = ".exe"
		}
	}

	runArgs := append([]interface{}{}, "./sitty"+exe)
	runArgs = append(runArgs, extraArgs...)

	Task("build").
		Exec("go", "build", "./cmd/sitty")

	Task("watch").
		Watch("cmd/*/*", "cmd/*/*/*", "pkg/*/*", "pkg/*/*/*", "pkg/*/*/*/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")

	Task("watch-only").
		Watch("cmd/*/*", "cmd/*/*/*", "pkg/*/*", "pkg/*/*/*", "pkg/*/*/*/*").
		Run("build")

	Task("run").
		Exec(runArgs...)

	Go()
}

package main

/*
	The application's business logic shall be implemented here in cli.actions.go in a function similar to `actionGreet` and `actionVersion` below.
	They must be added as a parameter in main.NewApplication where all possible application commands are defined.
*/

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
)

func actionGreet(c *cli.Context) error {
	if len(greetNameValue) > 0 {
		greetNameValue = ", " + greetNameValue
	}
	fmt.Printf("Hello%s\n", greetNameValue)
	if greetAskMeValue {
		fmt.Println("How are you?")
	}
	return nil
}

func actionVersion(c *cli.Context) error {
	fmt.Println(AppName + ":")
	fmt.Printf("  Version: %s\n", AppSemVersion)
	fmt.Printf("  Go version: %s\n", runtime.Version())
	fmt.Printf("  Git commit: %s\n", GitCommit)
	fmt.Printf("  Built: %s\n", AppBuildDate)
	fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Build type: %s\n", AppBuildType)
	return nil
}

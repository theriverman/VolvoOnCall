package main

/*
	The application's business logic shall be implemented here in cli.actions.go in a function similar to `actionGreet` and `actionVersion` below.
	They must be added as a parameter in main.NewApplication where all possible application commands are defined.
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

func actionRegister(c *cli.Context) error {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	defaultConfPath := filepath.Join(homeDirPath, ".voc.conf")
	Config.Username = username
	Config.Password = password
	return Config.WriteToFile(defaultConfPath)
}

func actionListCars(c *cli.Context) error {
	return nil
}

func actionStatus(c *cli.Context) error {
	return nil
}

func actionTrips(c *cli.Context) error {
	return nil
}

func actionLock(c *cli.Context) error {
	return nil
}

func actionUnlock(c *cli.Context) error {
	return nil
}

func startHeater(c *cli.Context) error {
	return nil
}

func stopHeater(c *cli.Context) error {
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

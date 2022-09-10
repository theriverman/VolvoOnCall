package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// application details injected at build time
var (
	AppName          string        = ""
	AppBuildDate     string        = ""
	AppBuildType     string        = ""
	AppSemVersion    string        = ""
	AppCopyrightText string        = ""
	GitCommit        string        = ""
	Config           Configuration = Configuration{}
)

// application behaviour
var appVerboseMode bool = false

// demo values for the template
// Refer to the documentation of urfave/cli at https://github.com/urfave/cli
var username, password string

// NewApplication is the primary entrypoint to our CLI application. the base logic shall be implemented here
func NewApplication() *cli.App {
	return &cli.App{
		Name:    AppName,
		Usage:   fmt.Sprintf("The purpose of %s is not explained here yet", AppName),
		Version: AppSemVersion,
		// application-level flags can be define below. these are applicable during the whole runtime
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Destination: &appVerboseMode,
				Usage:       fmt.Sprintf("Runs %s in verbose mode", AppName),
			},
		},
		Before: func(c *cli.Context) error {
			// CLI flags are processed at this point. Consider configuring your logging level
			if appVerboseMode {
				fmt.Println("Verbose mode has been enabled")
			}
			// Load configuration from $HOME (if exists)
			homeDirPath, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			defaultConfPath := filepath.Join(homeDirPath, ".voc.conf")
			if _, err := os.Stat(defaultConfPath); err == nil {
				if err = Config.LoadFromFile(defaultConfPath); err != nil {
					// add logging here
					return err
				}
				fmt.Printf("Config: %+v\n", Config)
			} else if errors.Is(err, os.ErrNotExist) {
				// add verbose logging here stating that the default config file does not exist
				fmt.Println("$HOME/.voc.conf was not found")
			}
			return nil
		},
		Commands: []*cli.Command{
			// list
			{
				Name:   "list",
				Usage:  "List all cars associated with your Volvo On Call account",
				Action: actionListCars,
				Flags:  []cli.Flag{},
			},
			// status
			{
				Name:   "status",
				Usage:  "Print a brief overview about the cars",
				Action: actionStatus,
				Flags:  []cli.Flag{},
			},
			// trips
			{
				Name:   "trips",
				Usage:  "Print a brief overview about the last trips",
				Action: actionTrips,
				Flags:  []cli.Flag{},
			},
			// owntracks

			// print

			// lock/unlock
			{
				Name:   "lock",
				Usage:  "Lock the car",
				Action: actionLock,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "unlock",
				Usage:  "Unlock the car",
				Action: actionUnlock,
				Flags:  []cli.Flag{},
			},

			// heater [start/stop]

			// engine [start/stop]

			// honk and blink

			// call (method)

			// mgtt
			{
				Name:   "register",
				Usage:  "Store your VOC username and password in a file at your $HOME folder",
				Action: actionRegister,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "username",
						Usage:       "Your Volvo On Call username",
						Value:       "",
						Destination: &username,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "password",
						Usage:       "Your Volvo On Call password",
						Value:       "",
						Destination: &password,
						Required:    true,
					},
				},
			},
			{
				Name:   "version",
				Usage:  fmt.Sprintf("Show the %s version information (detailed)", AppName),
				Action: actionVersion,
			},
		},
		Copyright: AppCopyrightText,
		// see the urfave/cli documentation for all possible options: https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	}
}

func main() {
	// if you're using the template as intended, this main() function shouldn't be modified at all
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Prints version information of go-socks5-cli and quit",
	}

	app := NewApplication()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

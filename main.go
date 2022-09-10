package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// application details injected at build time
var (
	AppName          string = ""
	AppBuildDate     string = ""
	AppBuildType     string = ""
	AppSemVersion    string = ""
	AppCopyrightText string = ""
	GitCommit        string = ""
)

// application behaviour
var appVerboseMode bool = false

// demo values for the template
// Refer to the documentation of urfave/cli at https://github.com/urfave/cli
var greetNameValue string
var greetAskMeValue bool

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
			return nil
		},
		Commands: []*cli.Command{
			// your custom application commands can be added here
			// see command `greet` below for a demonstration
			{
				Name:   "greet",
				Usage:  "Greets you appropriately to the configuration",
				Action: actionGreet,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "If you define this flag, then I will use this for greeting",
						Value:       "You",
						Destination: &greetNameValue,
						Required:    false,
					},
					&cli.BoolFlag{
						Name:        "ask-me",
						Usage:       "If you set this flag, then I will ask about your wellbeing",
						Value:       false,
						Destination: &greetAskMeValue,
						Required:    false,
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

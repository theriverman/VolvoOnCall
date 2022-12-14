package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	vocdriver "github.com/theriverman/VolvoOnCall"
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

// core handles
var client *vocdriver.Client

// application behaviour
var appVerboseMode bool = false

// runtime values
var selectedVin string = ""
var asJson bool = false
var customAttributes *cli.StringSlice = &cli.StringSlice{}

// NewApplication is the primary entrypoint to our CLI application. the base logic shall be implemented here
func NewApplication() *cli.App {
	return &cli.App{
		Name:    AppName,
		Usage:   "A CLI application to interact with Volvo Cars (On Call) services",
		Version: AppSemVersion,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Destination: &appVerboseMode,
				Usage:       fmt.Sprintf("Runs %s in verbose mode", AppName),
			},
			&cli.StringFlag{
				Name:        "username",
				Destination: &Config.Username,
				Usage:       "Volvo On Call username",
			},
			&cli.StringFlag{
				Name:        "password",
				Destination: &Config.Password,
				Usage:       "Volvo On Call password",
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
			} else if errors.Is(err, os.ErrNotExist) {
				// add verbose logging here stating that the default config file does not exist
				fmt.Println("$HOME/.voc.conf was not found")
			}
			client = &vocdriver.Client{
				ServiceRegion: Config.Region,
				BaseURL:       Config.URL,
			}
			if err = client.Initialise(); err != nil {
				return err
			}
			client.Authenticate(Config.Username, Config.Password)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "cars",
				Usage:  "List all cars associated with your Volvo On Call account",
				Action: actionCars,
			},

			// attributes
			{
				Name:   "attributes",
				Usage:  "Get all attributes of a select car",
				Action: actionAttributes,
				Before: selectVinOrThrowError,
				Flags: append(commonFlagsVin(), []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Usage:       "Return raw JSON response",
						Value:       false,
						Destination: &asJson,
					},
					&cli.StringSliceFlag{
						Name:        "attributes",
						Usage:       "Comma-separated JSON parameters to return",
						Destination: customAttributes,
					},
				}...),
			},

			// lock/unlock
			{
				Name:   "lock",
				Usage:  "Lock the car",
				Action: actionLock,
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
			},
			{
				Name:   "unlock",
				Usage:  "Unlock the car",
				Action: actionUnlock,
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
			},

			// heater [start/stop]
			{
				Name:  "heater",
				Usage: "Start/Stop the car's heater/climate",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "vin",
						Usage:       "Identify the car by its VIN",
						Value:       "",
						Destination: &selectedVin,
					},
				},
				Before: selectVinOrThrowError,
				Subcommands: []*cli.Command{
					{
						Name:   "start",
						Usage:  "Start the car's heater",
						Action: actionStartHeater,
					},
					{
						Name:   "stop",
						Usage:  "Stop  the car's heater",
						Action: actionStopHeater,
					},
				},
			},

			// engine
			{
				Name:   "engine",
				Usage:  "Start/Stop the car's engine",
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
				Subcommands: []*cli.Command{
					{
						Name:   "start",
						Usage:  "Start the car's engine",
						Action: actionStartEngine,
					},
					{
						Name:   "stop",
						Usage:  "Stop  the car's engine",
						Action: actionStopEngine,
					},
				},
			},

			// preclimatization
			{
				Name:   "preclimatization",
				Usage:  "Start/Stop preclimatization",
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
				Subcommands: []*cli.Command{
					{
						Name:   "start",
						Usage:  "Start preclimatization",
						Action: actionStartPreclimatization,
					},
					{
						Name:   "stop",
						Usage:  "Stop  preclimatization",
						Action: actionStopPreclimatization,
					},
				},
			},

			// honk and blink
			{
				Name:   "blink",
				Usage:  "Flash the turn signals",
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
				Action: actionBlink,
			},

			{
				Name:   "honk",
				Usage:  "Honk the horn",
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
				Action: actionHonk,
			},

			// position
			{
				Name:   "position",
				Usage:  "Get the car's last known position",
				Action: actionPosition,
				Before: selectVinOrThrowError,
				Flags: append(commonFlagsVin(), []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Usage:       "Return raw JSON response",
						Value:       false,
						Destination: &asJson,
					},
				}...),
			},

			// status
			{
				Name:   "status",
				Usage:  "Get a brief overview about a select car",
				Action: actionStatus,
				Before: selectVinOrThrowError,
				Flags: append(commonFlagsVin(), []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Usage:       "Return raw JSON response",
						Value:       false,
						Destination: &asJson,
					},
					&cli.StringSliceFlag{
						Name:        "attributes",
						Usage:       "Comma-separated JSON parameters to return",
						Destination: customAttributes,
					},
				}...),
			},

			// trips
			{
				Name:   "trips",
				Usage:  "Print a brief overview about the last trips",
				Action: actionTrips,
				Before: selectVinOrThrowError,
				Flags: append(commonFlagsVin(), []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Usage:       "Return raw JSON response",
						Value:       false,
						Destination: &asJson,
					},
				}...),
			},
			// owntracks

			// call (method)

			// mqtt

			// register
			{
				Name:   "register",
				Usage:  fmt.Sprintf("Save your %s username and password in $HOME/.voc.conf", AppName),
				Action: actionRegister,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "username",
						Usage:       "Your Volvo On Call username",
						Value:       Config.Username,
						Destination: &Config.Username,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "password",
						Usage:       "Your Volvo On Call password",
						Value:       Config.Password,
						Destination: &Config.Password,
						Required:    true,
					},
				},
			},
			// charging locations
			{
				Name:   "charging",
				Usage:  "View/Change charging locations",
				Flags:  commonFlagsVin(),
				Before: selectVinOrThrowError,
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List charging locations",
						Action: actionListChargingLocations,
					},
					{
						Name:      "get",
						Usage:     "Get a charging location by its ID",
						Action:    actionGetChargingLocationById,
						UsageText: "Pass a charging location's ID to get its details. Use the list command to find charging locations",
					},
					{
						Name:  "delay",
						Usage: "Enable/Disable delay charging and optionally update start/stop times",
						Subcommands: []*cli.Command{
							{
								Name:      "enable",
								Usage:     "Enable charging delay for a charging location",
								Action:    actionEnableDelayCharging,
								UsageText: "Pass a charging location's ID after the `enable` command. You can also pass a StartTime and a StopTime\nExample #1: voc charging delay enable 4075649\nExample #2: voc charging delay enable 4075649 22:35 06:50",
							},
							{
								Name:      "disable",
								Usage:     "Disable charging delay for a charging location",
								Action:    actionDisableDelayCharging,
								UsageText: "Pass a charging location's ID after the `disable` command\nFor example: voc charging delay disable 4075649",
							},
							{
								Name:      "update",
								Usage:     "Update both start time and stop time for a charging location",
								Action:    actionUpdateDelayCharging,
								UsageText: "Pass a charging location's ID after the `update` command along with a StartTime and a StopTime\nFor example: voc charging delay update 4075649 22:35 06:50",
							},
						},
					},
				},
			},

			// VOC version
			{
				Name:   "version",
				Usage:  fmt.Sprintf("Show the %s version information (detailed)", AppName),
				Action: actionVersion,
			},
		},
		Copyright: AppCopyrightText,
	}
}

func main() {
	// if you're using the template as intended, this main() function shouldn't be modified at all
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   fmt.Sprintf("Prints version information of %s and quit", AppName),
	}

	app := NewApplication()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

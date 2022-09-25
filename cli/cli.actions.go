package main

/*
	The application's business logic shall be implemented here in cli.actions.go in a function similar to `actionGreet` and `actionVersion` below.
	They must be added as a parameter in main.NewApplication where all possible application commands are defined.
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

func actionCars(c *cli.Context) error {
	account, err := client.CustomerAccount.GetAccount()
	if err != nil {
		return err
	}
	if err = account.RetrieveHyperlinks(); err != nil {
		return err
	}
	vehicles, err := account.GetVehicles()
	if err != nil {
		return err
	}

	fmt.Printf("Cars associated to Volvo Account(%s):\n", account.Username)
	fmt.Println("-----------------------------------" + strings.Repeat("-", len(account.Username)))
	for _, vehicle := range vehicles {
		if err = vehicle.RetrieveHyperlinks(); err != nil {
			return err
		}
		fmt.Printf("  * %s (%s)\n", vehicle.VehicleID, vehicle.Attributes.RegistrationNumber)
	}

	return nil
}

func actionRegister(c *cli.Context) error {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	defaultConfPath := filepath.Join(homeDirPath, ".voc.conf")
	return Config.WriteToFile(defaultConfPath)
}

func actionStatus(c *cli.Context) error {
	vehicle, err := client.Vehicles.GetVehicleByVIN(selectedVin)
	if err != nil {
		return err
	}
	s, err := json.MarshalIndent(vehicle.Status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(s))
	return nil
}

func actionTrips(c *cli.Context) error {
	trips, err := client.Vehicles.GetVehicleTripsByVIN(selectedVin)
	if err != nil {
		return err
	}
	s, err := json.MarshalIndent(trips.Trips, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(s))
	return nil
}

func actionLock(c *cli.Context) error {
	status, err := client.Vehicles.LockVehicle(selectedVin)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionUnlock(c *cli.Context) error {
	status, err := client.Vehicles.UnlockVehicle(selectedVin)
	if err != nil {
		return err
	}
	fmt.Println("Within 2 minutes press once gently on the rubberised pressure plate underneath the boot lid handle to unlock the car")
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionStartHeater(c *cli.Context) error {
	status, err := client.Vehicles.StartHeater(selectedVin)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionStopHeater(c *cli.Context) error {
	status, err := client.Vehicles.StopHeater(selectedVin)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionStartEngine(c *cli.Context) error {
	status, err := client.Vehicles.StartEngine(selectedVin)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionStopEngine(c *cli.Context) error {
	status, err := client.Vehicles.StopEngine(selectedVin)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionBlink(c *cli.Context) error {
	status, err := client.Vehicles.BlinkLights(selectedVin, nil)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionHonk(c *cli.Context) error {
	status, err := client.Vehicles.HonkAndBlink(selectedVin, nil)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
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

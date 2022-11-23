package main

/*
	The application's business logic shall be implemented here in cli.actions.go in a function similar to `actionGreet` and `actionVersion` below.
	They must be added as a parameter in main.NewApplication where all possible application commands are defined.
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/tabwriter"
	"unicode"

	vocdriver "github.com/theriverman/VolvoOnCall"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func actionAttributes(c *cli.Context) error {
	vehicle, err := client.Vehicles.GetVehicleByVIN(selectedVin)
	if err != nil {
		return err
	}
	if asJson {
		s, err := json.MarshalIndent(vehicle.Attributes, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(s))
		return nil
	}
	// for more advanced query options, see the Path Syntax at https://github.com/tidwall/gjson
	if len(customAttributes.Value()) > 0 {
		s, err := json.MarshalIndent(vehicle.Attributes, "", "\t")
		if err != nil {
			return err
		}
		for _, v := range customAttributes.Value() {
			fmt.Printf("%s: %s\n", v, gjson.GetBytes(s, v).String())
		}
		return nil
	}

	// default mode - print all attributes
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	fmt.Fprintf(w, "\n %s\t%s\t", "Car Attribute", "Attribute Value")
	fmt.Fprintf(w, "\n %s\t%s\t", "----", "----")
	fmt.Fprintf(w, "\n")
	v := reflect.ValueOf(*vehicle.Attributes)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if unicode.IsLower(rune(typeOfS.Field(i).Name[0])) {
			continue
		}
		fmt.Fprintf(w, "%s:\t %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	return nil
}

func actionPosition(c *cli.Context) error {
	pos, err := client.Vehicles.GetVehiclePositionByVIN(selectedVin)
	if err != nil {
		return err
	}
	if asJson {
		s, err := json.MarshalIndent(pos, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(s))
		return nil
	}

	// POSITION
	fmt.Println("Position")
	fmt.Printf("  - Longitude:\t%.15f\n", pos.Position.Longitude)
	fmt.Printf("  - Latitude:\t%.15f\n", pos.Position.Latitude)
	fmt.Printf("  - Timestamp:\t%s\n", pos.Position.Timestamp)
	// fmt.Printf("  - Speed:\t%s\n", (pos.Position.Speed))   // TODO: enable when actual type is known
	// fmt.Printf("  - Heading:\t%s\n", pos.Position.Heading) // TODO: enable when actual type is known
	fmt.Printf("  - Maps URL:\thttps://www.google.com/maps/place/%.15f,%.15f\n", pos.Position.Latitude, pos.Position.Longitude)

	// CALCULATED POSITION
	if pos.CalculatedPosition.Longitude > 1.0 {
		fmt.Println("\nCalculated Position")
		fmt.Printf("  - Longitude:\t%.15f\n", pos.CalculatedPosition.Longitude)
		fmt.Printf("  - Latitude:\t%.15f\n", pos.CalculatedPosition.Latitude)
		fmt.Printf("  - Timestamp:\t%s\n", pos.CalculatedPosition.Timestamp)
		// fmt.Printf("  - Speed:\t%s\n", pos.CalculatedPosition.Speed)     // TODO: enable when actual type is known
		// fmt.Printf("  - Heading:\t%s\n", pos.CalculatedPosition.Heading) // TODO: enable when actual type is known
		fmt.Printf("  - Maps URL:\thttps://www.google.com/maps/place/%.15f,%.15f\n", pos.CalculatedPosition.Latitude, pos.CalculatedPosition.Longitude)
	} else {
		fmt.Println("\nCalculated Position is not available")
	}

	return nil
}

func actionStatus(c *cli.Context) error {
	vehicle, err := client.Vehicles.GetVehicleByVIN(selectedVin)
	if err != nil {
		return err
	}
	if asJson {
		s, err := json.MarshalIndent(vehicle.Status, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(s))
		return nil
	}
	// for more advanced query options, see the Path Syntax at https://github.com/tidwall/gjson
	if len(customAttributes.Value()) > 0 {
		s, err := json.MarshalIndent(vehicle.Status, "", "\t")
		if err != nil {
			return err
		}
		for _, v := range customAttributes.Value() {
			fmt.Printf("%s: %s\n", v, gjson.GetBytes(s, v).String())
		}
		return nil
	}

	// default mode - print select attributes
	fmt.Printf("Average Fuel Consumption:\t%.1f l/100 km\n", vehicle.Status.AverageFuelConsumption/10)
	fmt.Printf("Average Speed:\t\t\t%d km/h\n", vehicle.Status.AverageSpeed)
	fmt.Printf("Brake Fluid:\t\t\t%s\n", vehicle.Status.BrakeFluid)
	if len(vehicle.Status.BulbFailures) > 0 {
		fmt.Println("Bulb Failures:")
		for _, failure := range vehicle.Status.BulbFailures {
			fmt.Printf("\t %s\n", failure)
		}
	}
	fmt.Printf("Car Locked:\t\t\t%t\n", vehicle.Status.CarLocked)
	fmt.Printf("Distance to Empty:\t\t%d km\n", vehicle.Status.DistanceToEmpty)
	doors := vehicle.Status.Doors
	if doors.HoodOpen || doors.FrontLeftDoorOpen || doors.FrontRightDoorOpen || doors.RearLeftDoorOpen || doors.RearRightDoorOpen || doors.TailgateOpen {
		fmt.Printf("Doors Open:\t\t\t%t\n", vehicle.Status.CarLocked)
		fmt.Printf("\t Hood Open: %t\n", doors.HoodOpen)
		fmt.Printf("\t Front Left Door Open: %t\n", doors.FrontLeftDoorOpen)
		fmt.Printf("\t Front Right Door Open: %t\n", doors.FrontRightDoorOpen)
		fmt.Printf("\t Rear Left Door Open: %t\n", doors.RearLeftDoorOpen)
		fmt.Printf("\t Rear Right Door Open: %t\n", doors.RearRightDoorOpen)
		fmt.Printf("\t Tailgate Open: %t\n", doors.TailgateOpen)
	} else {
		fmt.Println("Doors Open:\t\t\tNone")
	}
	fmt.Printf("Engine Running:\t\t\t%t\n", vehicle.Status.EngineRunning)
	fmt.Printf("Fuel Amount [l]:\t\t%d l\n", vehicle.Status.FuelAmount)
	fmt.Printf("Fuel Amount [%%]:\t\t%d%%\n", vehicle.Status.FuelAmountLevel)
	return actionPosition(c)
}

func actionTrips(c *cli.Context) error {
	trips, err := client.Vehicles.GetVehicleTripsByVIN(selectedVin)
	if err != nil {
		return err
	}
	if asJson {
		s, err := json.MarshalIndent(trips.Trips, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(s))
		return nil
	}

	// default mode
	p := message.NewPrinter(language.English)
	for i := len(trips.Trips) - 1; i >= 0; i-- { // ! reverse-loop !
		trip := trips.Trips[i]
		fmt.Printf("Trip %d\n", i+1)
		fmt.Printf("  ID: %d\n", trip.ID)
		fmt.Printf("  Name: %s\n", trip.Name)
		fmt.Printf("  Category: %s\n", trip.Category)
		fmt.Printf("  User Notes: %s\n", trip.UserNotes)
		for ii, tripDetail := range trip.TripDetails {
			fmt.Printf("  Trip Detail %d\n", ii+1)
			fmt.Printf("    Start Time: %s\n", tripDetail.StartTime)
			fmt.Printf("    End   Time: %s\n", tripDetail.EndTime)
			fmt.Printf("    Fuel Consumption: %.3f l\n", tripDetail.FuelConsumption/100)
			fmt.Printf("    Electrical Consumption: %.3f kWh\n", float64(tripDetail.ElectricalConsumption)/1000)
			fmt.Printf("    Electrical Regeneration: %.3f kWh\n", float64(tripDetail.ElectricalRegeneration)/1000)
			fmt.Printf("    Distance: %.3f km\n", tripDetail.Distance/1000)
			p.Printf("    Start Odometer: %d km\n", tripDetail.StartOdometer/1000)
			p.Printf("    End   Odometer: %d km\n", tripDetail.EndOdometer/1000)
			fmt.Printf("    Start Position: %+v\n", tripDetail.StartPosition)
			fmt.Printf("    End   Position: %+v\n", tripDetail.EndPosition)
			fmt.Println()
		}
		fmt.Println()
	}
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

func actionStartPreclimatization(c *cli.Context) error {
	status, err := client.Vehicles.StartPreclimatization(selectedVin)
	if err != nil {
		return err
	}
	return client.Vehicles.EvaluateServiceStatusAuto(status)
}

func actionStopPreclimatization(c *cli.Context) error {
	status, err := client.Vehicles.StopPreclimatization(selectedVin)
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

func actionListChargingLocations(c *cli.Context) error {
	chargingLocations, err := client.Vehicles.GetChargingLocations(selectedVin)
	if err != nil {
		return err
	}
	for _, cl := range chargingLocations.ChargingLocations {
		fmt.Printf("Name: %s\n", cl.Name)
		fmt.Printf("  ID: %s\n", path.Base(cl.ChargeLocation))
		fmt.Printf("  Status: %s\n", cl.Status)
		fmt.Printf("  Plug in reminder enabled: %t\n", cl.PlugInReminderEnabled)
		fmt.Printf("  Vehicle is at charging location: %t\n", cl.VehicleAtChargingLocation)
		fmt.Printf("  Position:\n")
		fmt.Printf("    - Street Address:\t%s\n", cl.Position.StreetAddress)
		fmt.Printf("    - City:\t\t%s\n", cl.Position.City)
		fmt.Printf("    - Postal Code:\t%s\n", cl.Position.PostalCode)
		fmt.Printf("    - Region:\t\t%s\n", cl.Position.Region)
		fmt.Printf("    - Country Code:\t%s\n", cl.Position.ISO2CountryCode)
		fmt.Printf("    - Longitude:\t%.15f\n", cl.Position.Longitude)
		fmt.Printf("    - Latitude:\t\t%.15f\n", cl.Position.Latitude)
		fmt.Printf("    - Maps URL:\t\thttps://www.google.com/maps/place/%.15f,%.15f\n", cl.Position.Latitude, cl.Position.Longitude)
		fmt.Printf("---------------------------------------------------------------------------------------\n\n")
	}
	return nil
}

func actionGetChargingLocationById(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("you must provide a charging location id. see --help for more details")
	}
	cl, err := client.Vehicles.GetChargingLocation(selectedVin, c.Args().First())
	if err != nil {
		return err
	}
	fmt.Printf("Name: %s\n", cl.Name)
	fmt.Printf("  ID: %s\n", path.Base(cl.ChargeLocation))
	fmt.Printf("  Status: %s\n", cl.Status)
	fmt.Printf("  Plug in reminder enabled: %t\n", cl.PlugInReminderEnabled)
	fmt.Printf("  Vehicle is at charging location: %t\n", cl.VehicleAtChargingLocation)
	fmt.Printf("  Position:\n")
	fmt.Printf("    - Street Address:\t%s\n", cl.Position.StreetAddress)
	fmt.Printf("    - City:\t\t%s\n", cl.Position.City)
	fmt.Printf("    - Postal Code:\t%s\n", cl.Position.PostalCode)
	fmt.Printf("    - Region:\t\t%s\n", cl.Position.Region)
	fmt.Printf("    - Country Code:\t%s\n", cl.Position.ISO2CountryCode)
	fmt.Printf("    - Longitude:\t%.15f\n", cl.Position.Longitude)
	fmt.Printf("    - Latitude:\t\t%.15f\n", cl.Position.Latitude)
	fmt.Printf("    - Maps URL:\t\thttps://www.google.com/maps/place/%.15f,%.15f\n", cl.Position.Latitude, cl.Position.Longitude)
	fmt.Printf("  ---------------------------------------------------------------------------------------------\n\n")
	return nil
}

func actionEnableDelayCharging(c *cli.Context) error {
	var (
		chargingId, startTime, stopTime string
	)
	dc := vocdriver.DelayCharging{
		Enabled: true,
	}
	switch c.Args().Len() {
	case 0:
		return fmt.Errorf("you must provide a charging location id. see --help for more details")
	case 1:
		chargingId = c.Args().First()
		cl, err := client.Vehicles.GetChargingLocation(selectedVin, chargingId)
		if err != nil {
			return err
		}
		dc.StartTime = cl.DelayCharging.StartTime
		dc.StopTime = cl.DelayCharging.StopTime
	case 2:
		return fmt.Errorf("unexpected number of arguments were passed. minimum 1 or exactly 3 allowed")
	case 3:
		chargingId = c.Args().First()
		startTime = c.Args().Get(1)
		stopTime = c.Args().Get(2)
		dc.StartTime = startTime
		dc.StopTime = stopTime
	default:
		return fmt.Errorf("unexpected number of arguments were passed. minimum 1 or exactly 3 allowed")
	}

	vehicle, err := client.Vehicles.GetVehicleByVIN(selectedVin)
	if err != nil {
		return err
	}

	if _, err = vehicle.SetDelayCharging(chargingId, &dc); err != nil {
		return err
	}
	return nil
}

func actionDisableDelayCharging(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("you must provide a charging location id. see --help for more details")
	}
	dc := vocdriver.DelayCharging{
		Enabled: false,
	}
	vehicle, err := client.Vehicles.GetVehicleByVIN(selectedVin)
	if err != nil {
		return err
	}
	if _, err = vehicle.SetDelayCharging(c.Args().First(), &dc); err != nil {
		return err
	}
	return nil
}

func actionUpdateDelayCharging(c *cli.Context) error {
	var (
		chargingId, startTime, stopTime string
	)
	dc := vocdriver.DelayCharging{}
	switch c.Args().Len() {
	case 3:
		chargingId = c.Args().First()
		startTime = c.Args().Get(1)
		stopTime = c.Args().Get(2)

		cl, err := client.Vehicles.GetChargingLocation(selectedVin, chargingId)
		if err != nil {
			return err
		}
		dc.Enabled = cl.DelayCharging.Enabled // keep current status
		dc.StartTime = startTime
		dc.StopTime = stopTime
	default:
		return fmt.Errorf("you must provide: charging location id + start time + stop time. see --help for more details")
	}

	vehicle, err := client.Vehicles.GetVehicleByVIN(selectedVin)
	if err != nil {
		return err
	}

	if _, err = vehicle.SetDelayCharging(chargingId, &dc); err != nil {
		return err
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

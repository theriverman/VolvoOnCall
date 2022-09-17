package vocdriver

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestClient_Initialise(t *testing.T) {
	err := godotenv.Load("/workspaces/VolvoOnCall/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("username")
	password := os.Getenv("password")

	client := Client{}
	client.Initialise()
	client.Authenticate(username, password)

	accounts, err := client.CustomerAccounts.GetAccount()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	t.Logf("Response: %+v\n", accounts)

	accountVehicleRelations := []VehicleAccountRelation{}
	for _, relationId := range accounts.GetAccountVehicleRelations() {
		vehicleAccRel, err := client.VehicleAccountRelation.GetById(relationId)
		if err != nil {
			t.Errorf("%v\n", err)
			continue
		}
		accountVehicleRelations = append(accountVehicleRelations, *vehicleAccRel)
	}

	vehicles := []Vehicle{}
	for _, avr := range accountVehicleRelations {
		vehicle, err := client.Vehicles.GetVehicleByVIN(avr.VehicleID)
		if err != nil {
			t.Errorf("%v\n", err)
			continue
		}
		vehicles = append(vehicles, *vehicle)
	}

	t.Logf("My Vehicles:\n")
	for _, vehicle := range vehicles {
		attributes, err := vehicle.GetAttributes()
		if err != nil {
			t.Errorf("%v\n", err)
			continue
		}
		t.Logf("  * %s (%s)\n", vehicle.VehicleID, attributes.RegistrationNumber)
	}
}

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

	client, err := NewClient(username, password)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	account, err := client.CustomerAccount.GetAccount()
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	vehicles, err := account.GetVehicles()
	if err != nil {
		t.Fatalf("%v\n", err)
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

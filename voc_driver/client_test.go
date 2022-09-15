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

	c := Client{}
	c.Initialise()
	c.Authenticate(username, password)

	accounts, err := c.CustomerAccounts.Get()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	t.Logf("Response: %+v\n", accounts)
}

# VolvoOnCall
Library and [CLI application](./voc/) written in Go to interact with the Volvo Cars (Volvo On Call) services.

This project was inspired by [molobrakos/volvooncall](https://github.com/molobrakos/volvooncall), and it aims to maintain a certain level of compatibility with it both API and configuration wise.

For more about the CLI application, see [cli/README.md](./voc/README.md).

# Get the Library
```bash
go get github.com/theriverman/VolvoOnCall
```

# Getting Started
```go
import (
  "fmt"
  vocdriver "github.com/theriverman/VolvoOnCall"
)

client, err := NewClient("your-volvo-on-call-username", "your password")
if err != nil {
  fmt.Printf("%v\n", err)
}

account, err := client.CustomerAccount.GetAccount()
if err != nil {
  fmt.Printf("%v\n", err)
}

vehicles, err := account.GetVehicles()
if err != nil {
  fmt.Printf("%v\n", err)
}

fmt.Logf("My Vehicles:\n")
for _, vehicle := range vehicles {
  if err = vehicle.RetrieveHyperlinks(); err != nil {
    fmt.Printf("%v\n", err)
    continue
  }
  fmt.Printf("  * %s (%s)\n", vehicle.VehicleID, vehicle.Attributes.RegistrationNumber)
  fmt.Printf("    - IsHeaterSupported: %t\n", vehicle.IsHeaterSupported())
  status, err := vehicle.BlinkLights(nil)
  if err != nil {
    fmt.Panicf("%v", err)
  }
  if err = client.Vehicles.EvaluateServiceStatusAuto(status); err != nil {
    fmt.Panicf("%v", err)
  }
}
```

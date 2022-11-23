package vocdriver

import (
	"fmt"
	"log"
	"time"
)

var ServiceTypeMap map[string]string = map[string]string{
	"RDL":   "Lock Vehicle",
	"RDU":   "Unlock Vehicle",
	"RHBLF": "Blink Lights",
}

type VehiclesService struct {
	client   *Client
	Endpoint string
}

/*
	Low-level Functions
*/

func (v *VehiclesService) GetVehicleByVIN(vin string) (vehicle *Vehicle, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin)
	if _, err = v.client.Request.Get(url, &vehicle); err != nil {
		return nil, err
	}
	vehicle.client = v.client
	err = vehicle.RetrieveHyperlinks()
	return
}

func (v *VehiclesService) GetVehicleByHyperlink(url string) (vehicle *Vehicle, err error) {
	if url == "" {
		return nil, fmt.Errorf("url must not be empty")
	}
	if _, err = v.client.Request.Get(url, &vehicle); err != nil {
		return nil, err
	}
	vehicle.client = v.client
	err = vehicle.RetrieveHyperlinks()
	return
}

func (v *VehiclesService) GetVehicleAttributesByVIN(vin string) (attributes *VehicleAttributes, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "attributes")
	if _, err = v.client.Request.Get(url, &attributes); err != nil {
		return nil, err
	}
	attributes.client = v.client
	return
}

func (v *VehiclesService) GetVehicleStatusByVIN(vin string) (status *VehicleStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "status")
	if _, err = v.client.Request.Get(url, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) GetVehiclePositionByVIN(vin string) (position *VehiclePosition, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "position")
	if _, err = v.client.Request.Get(url, &position); err != nil {
		return nil, err
	}
	position.client = v.client
	return
}

func (v *VehiclesService) GetVehicleTripsByVIN(vin string) (trips *VehicleTrips, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "trips")
	if _, err = v.client.Request.Get(url, &trips); err != nil {
		return nil, err
	}
	trips.client = v.client
	return
}

// GetServiceStatus retrieves the current status of an async operation (typically an action sent to a vehicle)
func (v *VehiclesService) GetServiceStatus(url string) (vss *VehicleServiceStatus, err error) {
	if _, err = v.client.Request.Get(url, &vss); err != nil {
		return nil, err
	}
	vss.client = v.client
	return
}

// EvaluateServiceStatusAuto hangs the main application waiting for the requested operation to finish
// During this time the Service API is polled every second and the response is evaluated
//   - if the request timeouts (default: 30s), an error is returned
//   - if the request fails, an error is returned
func (v *VehiclesService) EvaluateServiceStatusAuto(vss *VehicleServiceStatus) (err error) {
	timeoutSeconds := 30
	if ServiceTypeMap[vss.ServiceType] == "Unlock Vehicle" {
		vehicle, err := v.client.Vehicles.GetVehicleByVIN(vss.VehicleID)
		if err != nil {
			return fmt.Errorf("failed to retrieve vehicle details for %s", vss.VehicleID)
		}
		timeoutSeconds = vehicle.Attributes.UnlockTimeFrame
		log.Printf("value of timeoutSeconds increased to %d to match the vehicle's unlockTimeFrame value", timeoutSeconds)
	}
	return v.EvaluateServiceStatus(vss, timeoutSeconds)
}

func (v *VehiclesService) EvaluateServiceStatus(vss *VehicleServiceStatus, timeoutSeconds int) (err error) {
	c := 0
	for {
		if c == timeoutSeconds {
			return fmt.Errorf("request timeout (%ds)", timeoutSeconds)
		}
		if c > 0 {
			if err = vss.Refresh(); err != nil {
				return
			}
		}
		switch vss.Status {
		case "Started":
			time.Sleep(1 * time.Second)
			c++
			continue
		case "MessageDelivered":
			time.Sleep(1 * time.Second)
			c++
			continue
		case "Successful":
			return nil
		case "Failed":
			return fmt.Errorf("request (%s) failed: %s", ServiceTypeMap[vss.ServiceType], vss.FailureReason)
		default:
			return fmt.Errorf("request (%s) failed with status (%s): %s", ServiceTypeMap[vss.ServiceType], vss.Status, vss.FailureReason)
		}
	}
}

// BlinkLights blinks the lights on the car without sounding the horn
//
// This API requires sending the client's position. You have two options:
//   - Share your position by passing a valid *Position struct.
//   - Pass in `nil` and the actual own position of the car will be sent used
func (v *VehiclesService) BlinkLights(vin string, position *Position) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "honk_blink", "lights")
	if position == nil {
		position = &Position{}
		vehiclePosition, err := v.GetVehiclePositionByVIN(vin)
		if err != nil {
			return nil, err
		}
		position.Longitude = vehiclePosition.Position.Longitude
		position.Latitude = vehiclePosition.Position.Latitude
	}
	payload := map[string]float64{
		"clientAccuracy":  0.0,
		"clientLatitude":  position.Latitude,
		"clientLongitude": position.Longitude,
	}
	if _, err = v.client.Request.Post(url, payload, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

// HonkAndBlink honks and blinks the car
//
// This API requires sending the client's position. You have two options:
//   - Share your position by passing a valid *Position struct.
//   - Pass in `nil` and the actual own position of the car will be sent used
func (v *VehiclesService) HonkAndBlink(vin string, position *Position) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "honkAndBlink")
	if position == nil {
		position = &Position{}
		vehiclePosition, err := v.GetVehiclePositionByVIN(vin)
		if err != nil {
			return nil, err
		}
		position.Longitude = vehiclePosition.Position.Longitude
		position.Latitude = vehiclePosition.Position.Latitude
	}
	payload := map[string]float64{
		"clientAccuracy":  0.0,
		"clientLatitude":  position.Latitude,
		"clientLongitude": position.Longitude,
	}
	if _, err = v.client.Request.Post(url, payload, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) LockVehicle(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "lock")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) UnlockVehicle(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "unlock")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) StartEngine(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "engine", "start")
	if _, err = v.client.Request.Post(url, map[string]int{"runtime": 15}, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) StopEngine(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "engine", "stop")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) StartHeater(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "heater", "start")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) StopHeater(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "heater", "stop")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) StartPreclimatization(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "preclimatization", "start")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) StopPreclimatization(vin string) (status *VehicleServiceStatus, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "preclimatization", "stop")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

/*
Charge Locations
*/
func (v *VehiclesService) GetChargingLocations(vin string) (chargingLocations *ChargingLocations, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "chargeLocations")
	if _, err = v.client.Request.Get(url, &chargingLocations); err != nil {
		return nil, err
	}
	chargingLocations.client = v.client
	for i := 0; i < len(chargingLocations.ChargingLocations); i++ {
		chargingLocations.ChargingLocations[i].client = v.client
	}
	return
}

func (v *VehiclesService) GetChargingLocation(vin, chargingId string) (chargingLocation *ChargingLocation, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	if chargingId == "" {
		return nil, fmt.Errorf("chargingId must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "chargeLocations", chargingId)
	if _, err = v.client.Request.Get(url, &chargingLocation); err != nil {
		return nil, err
	}
	chargingLocation.client = v.client
	return
}

func (v *VehiclesService) UpdateChargingLocation(vin, chargingId string, chargingLocation *ChargingLocation) (chargingLocationResponse *ChargingLocation, err error) {
	if vin == "" {
		return nil, fmt.Errorf("vin must not be empty")
	}
	if chargingId == "" {
		return nil, fmt.Errorf("chargingId must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "chargeLocations", chargingId)
	if _, err = v.client.Request.Put(url, &chargingLocation, &chargingLocationResponse); err != nil {
		return nil, err
	}
	chargingLocationResponse.client = v.client
	return
}

/*
	Utils
*/

func (v *VehiclesService) RetrieveServiceStatus(vin, customerServiceId string) (status *VehicleServiceStatus, err error) {
	if vin == "" || customerServiceId == "" {
		return nil, fmt.Errorf("vin and customerServiceId must not be empty")
	}
	url := v.client.MakeURL(v.Endpoint, vin, "services", customerServiceId)
	if _, err = v.client.Request.Get(url, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

type Vehicle struct {
	Attributes                       *VehicleAttributes
	Status                           *VehicleStatus
	VehicleAccountRelations          []AccountVehicleRelation
	HyperlinkAttributes              string   `json:"attributes"`              // url
	HyperlinkStatus                  string   `json:"status"`                  // url
	HyperlinkVehicleAccountRelations []string `json:"vehicleAccountRelations"` // url
	VehicleID                        string   `json:"vehicleId"`               // vin
	vehicleAccountRelationsRetrieved bool
	client                           *Client // added for interface simplification
}

func (v *Vehicle) RetrieveHyperlinks() (err error) {
	if v.Attributes == nil {
		if v.Attributes, err = v.client.Vehicles.GetVehicleAttributesByVIN(v.VehicleID); err != nil {
			return
		}
	}
	if v.Status == nil {
		if v.Status, err = v.client.Vehicles.GetVehicleStatusByVIN(v.VehicleID); err != nil {
			return
		}
	}
	if !v.vehicleAccountRelationsRetrieved {
		for _, url := range v.HyperlinkVehicleAccountRelations {
			relation, err := v.client.AccountVehicleRelation.GetByHyperlink(url)
			if err != nil {
				return err
			}
			v.VehicleAccountRelations = append(v.VehicleAccountRelations, *relation)
		}
		v.vehicleAccountRelationsRetrieved = true
	}
	return
}

/*
Actions/Operations
*/
func (v Vehicle) BlinkLights(position *Position) (status *VehicleServiceStatus, err error) {
	return v.client.Vehicles.BlinkLights(v.VehicleID, position)
}

func (v *Vehicle) GetPosition() (position *VehiclePosition, err error) {
	return v.client.Vehicles.GetVehiclePositionByVIN(v.VehicleID)
}

func (v *Vehicle) GetTrips() (trips *VehicleTrips, err error) {
	return v.client.Vehicles.GetVehicleTripsByVIN(v.VehicleID)
}

func (v *Vehicle) Lock() (status *VehicleServiceStatus, err error) {
	if !v.IsLockSupported() {
		return nil, fmt.Errorf("lock/unlock is not supported by %s [%s]", v.Attributes.RegistrationNumber, v.Attributes.Vin)
	}
	return v.client.Vehicles.LockVehicle(v.VehicleID)
}

func (v Vehicle) UnlockVehicle() (status *VehicleServiceStatus, err error) {
	if !v.IsUnlockSupported() {
		return nil, fmt.Errorf("lock/unlock is not supported by %s [%s]", v.Attributes.RegistrationNumber, v.Attributes.Vin)
	}
	return v.client.Vehicles.UnlockVehicle(v.VehicleID)
}

func (v Vehicle) StartEngine() (status *VehicleServiceStatus, err error) {
	if !v.IsEngineStartSupported() {
		return nil, fmt.Errorf("engine start/stop is not supported by %s [%s]", v.Attributes.RegistrationNumber, v.Attributes.Vin)
	}
	return v.client.Vehicles.StartEngine(v.VehicleID)
}

func (v Vehicle) StopEngine() (status *VehicleServiceStatus, err error) {
	if !v.IsEngineStartSupported() {
		return nil, fmt.Errorf("engine start/stop is not supported by %s [%s]", v.Attributes.RegistrationNumber, v.Attributes.Vin)
	}
	return v.client.Vehicles.StopEngine(v.VehicleID)
}

func (v Vehicle) StartHeater() (status *VehicleServiceStatus, err error) {
	switch {
	case v.IsHeaterSupported():
		return v.client.Vehicles.StartHeater(v.VehicleID)
	case v.IsPreclimatizationSupported():
		return v.client.Vehicles.StartPreclimatization(v.VehicleID)
	default:
		return nil, fmt.Errorf("heater is not supported by %s [%s]", v.Attributes.RegistrationNumber, v.Attributes.Vin)
	}
}

func (v Vehicle) StopHeater() (status *VehicleServiceStatus, err error) {
	switch {
	case v.IsHeaterSupported():
		return v.client.Vehicles.StopHeater(v.VehicleID)
	case v.IsPreclimatizationSupported():
		return v.client.Vehicles.StopPreclimatization(v.VehicleID)
	default:
		return nil, fmt.Errorf("heater is not supported by %s [%s]", v.Attributes.RegistrationNumber, v.Attributes.Vin)
	}
}

func (v Vehicle) SetDelayCharging(chargingId string, delayCharging *DelayCharging) (chargingLocation *ChargingLocation, err error) {
	cl := ChargingLocation{
		Status:        "Accepted",
		DelayCharging: delayCharging,
	}
	return v.client.Vehicles.UpdateChargingLocation(v.VehicleID, chargingId, &cl)
}

/*
	Properties
*/

// IsHeaterSupported returns true if either Remote Heater or Preclimatization is supported
func (v Vehicle) IsHeaterSupported() bool {
	return v.Attributes.RemoteHeaterSupported || v.Attributes.PreclimatizationSupported
}

func (v Vehicle) IsHeaterOn() bool {
	switch v.Status.Heater.Status {
	case "on":
		return true
	case "off":
		return false
	default:
		log.Printf("IsHeaterOn :: unexpected response: %s\n", v.Status.Heater.Status)
		return false
	}
}

func (v Vehicle) IsLockSupported() bool {
	return v.Attributes.LockSupported
}

func (v Vehicle) IsUnlockSupported() bool {
	return v.Attributes.UnlockSupported
}

func (v Vehicle) IsLocked() bool {
	return v.Status.CarLocked
}

func (v Vehicle) IsRemoteHeaterSupported() bool {
	return v.Attributes.RemoteHeaterSupported
}

func (v Vehicle) IsPreclimatizationSupported() bool {
	return v.Attributes.PreclimatizationSupported
}

func (v Vehicle) IsJournalLogSupported() bool {
	return v.Attributes.JournalLogSupported
}

func (v Vehicle) IsJournalLogEnabled() bool {
	return v.Attributes.JournalLogEnabled
}

func (v Vehicle) IsHonkAndBlinkSupported() bool {
	return v.Attributes.HonkAndBlinkSupported
}

func (v Vehicle) IsEngineStartSupported() bool {
	return v.Attributes.EngineStartSupported
}

type VehicleAttributes struct {
	EngineCode                             string   `json:"engineCode"`
	ExteriorCode                           string   `json:"exteriorCode"`
	InteriorCode                           string   `json:"interiorCode"`
	TyreDimensionCode                      string   `json:"tyreDimensionCode"`
	TyreInflationPressureLightCode         string   `json:"tyreInflationPressureLightCode"`
	TyreInflationPressureHeavyCode         string   `json:"tyreInflationPressureHeavyCode"`
	GearboxCode                            string   `json:"gearboxCode"`
	FuelType                               string   `json:"fuelType"`
	FuelTankVolume                         int      `json:"fuelTankVolume"`
	GrossWeight                            int      `json:"grossWeight"`
	ModelYear                              int      `json:"modelYear"`
	VehicleType                            string   `json:"vehicleType"`
	VehicleTypeCode                        string   `json:"vehicleTypeCode"`
	NumberOfDoors                          int      `json:"numberOfDoors"`
	RegistrationNumber                     string   `json:"registrationNumber"`
	CarLocatorDistance                     int      `json:"carLocatorDistance"`
	HonkAndBlinkDistance                   int      `json:"honkAndBlinkDistance"`
	BCallAssistanceNumber                  string   `json:"bCallAssistanceNumber"`
	CarLocatorSupported                    bool     `json:"carLocatorSupported"`
	HonkAndBlinkSupported                  bool     `json:"honkAndBlinkSupported"`
	HonkAndBlinkVersionsSupported          []string `json:"honkAndBlinkVersionsSupported"`
	RemoteHeaterSupported                  bool     `json:"remoteHeaterSupported"`
	UnlockSupported                        bool     `json:"unlockSupported"`
	LockSupported                          bool     `json:"lockSupported"`
	JournalLogSupported                    bool     `json:"journalLogSupported"`
	AssistanceCallSupported                bool     `json:"assistanceCallSupported"`
	UnlockTimeFrame                        int      `json:"unlockTimeFrame"`
	VerificationTimeFrame                  int      `json:"verificationTimeFrame"`
	TimeFullyAccessible                    int      `json:"timeFullyAccessible"`
	TimePartiallyAccessible                int      `json:"timePartiallyAccessible"`
	SubscriptionType                       string   `json:"subscriptionType"`
	SubscriptionStartDate                  string   `json:"subscriptionStartDate"`
	SubscriptionEndDate                    string   `json:"subscriptionEndDate"`
	ServerVersion                          string   `json:"serverVersion"`
	Vin                                    string   `json:"VIN"`
	JournalLogEnabled                      bool     `json:"journalLogEnabled"`
	HighVoltageBatterySupported            bool     `json:"highVoltageBatterySupported"`
	MaxActiveDelayChargingLocations        int      `json:"maxActiveDelayChargingLocations"`
	PreclimatizationSupported              bool     `json:"preclimatizationSupported"`
	SendPOIToVehicleVersionsSupported      []string `json:"sendPOIToVehicleVersionsSupported"`
	ClimatizationCalendarVersionsSupported []string `json:"climatizationCalendarVersionsSupported"`
	ClimatizationCalendarMaxTimers         int      `json:"climatizationCalendarMaxTimers"`
	VehiclePlatform                        string   `json:"vehiclePlatform"`
	VinLower                               string   `json:"vin"`
	OverrideDelayChargingSupported         bool     `json:"overrideDelayChargingSupported"`
	EngineStartSupported                   bool     `json:"engineStartSupported"`
	StatusParkedIndoorSupported            bool     `json:"status.parkedIndoor.supported"`
	Country                                struct {
		Iso2 string `json:"iso2"`
	} `json:"country"`
	client *Client // added for interface simplification
}

func (va VehicleAttributes) VIN() string {
	switch {
	case len(va.Vin) > 0:
		return va.Vin
	case len(va.VinLower) > 0:
		return va.VinLower
	default:
		log.Panicln("vin could not be retrieved from VehicleAttributes")
		return ""
	}
}

type VehicleStatus struct {
	AverageFuelConsumption          float64  `json:"averageFuelConsumption"`
	AverageFuelConsumptionTimestamp string   `json:"averageFuelConsumptionTimestamp"`
	AverageSpeed                    int      `json:"averageSpeed"`
	AverageSpeedTimestamp           string   `json:"averageSpeedTimestamp"`
	BrakeFluid                      string   `json:"brakeFluid"`
	BrakeFluidTimestamp             string   `json:"brakeFluidTimestamp"`
	BulbFailures                    []string `json:"bulbFailures"`
	BulbFailuresTimestamp           string   `json:"bulbFailuresTimestamp"`
	CarLocked                       bool     `json:"carLocked"`
	CarLockedTimestamp              string   `json:"carLockedTimestamp"`
	ConnectionStatus                string   `json:"connectionStatus"`
	ConnectionStatusTimestamp       string   `json:"connectionStatusTimestamp"`
	DistanceToEmpty                 int      `json:"distanceToEmpty"`
	DistanceToEmptyTimestamp        string   `json:"distanceToEmptyTimestamp"`
	Doors                           struct {
		TailgateOpen       bool   `json:"tailgateOpen"`
		RearRightDoorOpen  bool   `json:"rearRightDoorOpen"`
		RearLeftDoorOpen   bool   `json:"rearLeftDoorOpen"`
		FrontRightDoorOpen bool   `json:"frontRightDoorOpen"`
		FrontLeftDoorOpen  bool   `json:"frontLeftDoorOpen"`
		HoodOpen           bool   `json:"hoodOpen"`
		Timestamp          string `json:"timestamp"`
	} `json:"doors"`
	EngineRunning            bool   `json:"engineRunning"`
	EngineRunningTimestamp   string `json:"engineRunningTimestamp"`
	FuelAmount               int    `json:"fuelAmount"`
	FuelAmountLevel          int    `json:"fuelAmountLevel"`
	FuelAmountLevelTimestamp string `json:"fuelAmountLevelTimestamp"`
	FuelAmountTimestamp      string `json:"fuelAmountTimestamp"`
	Heater                   struct {
		SeatSelection struct {
			FrontDriverSide    bool `json:"frontDriverSide"`
			FrontPassengerSide bool `json:"frontPassengerSide"`
			RearDriverSide     bool `json:"rearDriverSide"`
			RearPassengerSide  bool `json:"rearPassengerSide"`
			RearMid            bool `json:"rearMid"`
		} `json:"seatSelection"`
		Status string `json:"status"`
		Timer1 struct {
			Time  string `json:"time"`
			State bool   `json:"state"`
		} `json:"timer1"`
		Timer2 struct {
			Time  string `json:"time"`
			State bool   `json:"state"`
		} `json:"timer2"`
		Timestamp string `json:"timestamp"`
	} `json:"heater"`
	HvBattery struct {
		HvBatteryChargeStatusDerived          string `json:"hvBatteryChargeStatusDerived"`
		HvBatteryChargeStatusDerivedTimestamp string `json:"hvBatteryChargeStatusDerivedTimestamp"`
		HvBatteryChargeModeStatus             string `json:"hvBatteryChargeModeStatus"`
		HvBatteryChargeModeStatusTimestamp    string `json:"hvBatteryChargeModeStatusTimestamp"`
		HvBatteryChargeStatus                 string `json:"hvBatteryChargeStatus"`
		HvBatteryChargeStatusTimestamp        string `json:"hvBatteryChargeStatusTimestamp"`
		HvBatteryLevel                        int    `json:"hvBatteryLevel"`
		HvBatteryLevelTimestamp               string `json:"hvBatteryLevelTimestamp"`
		DistanceToHVBatteryEmpty              int    `json:"distanceToHVBatteryEmpty"`
		DistanceToHVBatteryEmptyTimestamp     string `json:"distanceToHVBatteryEmptyTimestamp"`
		HvBatteryChargeWarning                string `json:"hvBatteryChargeWarning"`
		HvBatteryChargeWarningTimestamp       string `json:"hvBatteryChargeWarningTimestamp"`
		TimeToHVBatteryFullyCharged           int    `json:"timeToHVBatteryFullyCharged"`
		TimeToHVBatteryFullyChargedTimestamp  string `json:"timeToHVBatteryFullyChargedTimestamp"`
	} `json:"hvBattery"`
	Odometer                           int    `json:"odometer"`
	OdometerTimestamp                  string `json:"odometerTimestamp"`
	ParkedIndoor                       bool   `json:"parkedIndoor"`
	ParkedIndoorTimestamp              string `json:"parkedIndoorTimestamp"`
	RemoteClimatizationStatus          string `json:"remoteClimatizationStatus"`
	RemoteClimatizationStatusTimestamp string `json:"remoteClimatizationStatusTimestamp"`
	ServiceWarningStatus               string `json:"serviceWarningStatus"`
	ServiceWarningStatusTimestamp      string `json:"serviceWarningStatusTimestamp"`
	TheftAlarm                         struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Timestamp string  `json:"timestamp"`
	} `json:"theftAlarm"`
	TimeFullyAccessibleUntil     string `json:"timeFullyAccessibleUntil"`
	TimePartiallyAccessibleUntil string `json:"timePartiallyAccessibleUntil"`
	TripMeter1                   int    `json:"tripMeter1"`
	TripMeter1Timestamp          string `json:"tripMeter1Timestamp"`
	TripMeter2                   int    `json:"tripMeter2"`
	TripMeter2Timestamp          string `json:"tripMeter2Timestamp"`
	WasherFluidLevel             string `json:"washerFluidLevel"`
	WasherFluidLevelTimestamp    string `json:"washerFluidLevelTimestamp"`
	Windows                      struct {
		FrontLeftWindowOpen  bool   `json:"frontLeftWindowOpen"`
		FrontRightWindowOpen bool   `json:"frontRightWindowOpen"`
		Timestamp            string `json:"timestamp"`
		RearLeftWindowOpen   bool   `json:"rearLeftWindowOpen"`
		RearRightWindowOpen  bool   `json:"rearRightWindowOpen"`
	} `json:"windows"`
	client *Client // added for interface simplification
}

type VehiclePosition struct {
	Position           Position `json:"position"`
	CalculatedPosition Position `json:"calculatedPosition"`
	client             *Client
}

type VehicleTrips struct {
	Trips []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Category     string `json:"category"`
		UserNotes    string `json:"userNotes"`
		Trip         string `json:"trip"`
		RouteDetails struct {
			Route          string `json:"route"`
			TotalWaypoints int    `json:"totalWaypoints"`
			BoundingBox    struct {
				MinLongitude float64 `json:"minLongitude"`
				MinLatitude  float64 `json:"minLatitude"`
				MaxLongitude float64 `json:"maxLongitude"`
				MaxLatitude  float64 `json:"maxLatitude"`
			} `json:"boundingBox"`
		} `json:"routeDetails,omitempty"`
		TripDetails []struct {
			FuelConsumption        float64 `json:"fuelConsumption"`
			ElectricalConsumption  float64 `json:"electricalConsumption"`
			ElectricalRegeneration float64 `json:"electricalRegeneration"`
			Distance               float64 `json:"distance"`
			StartOdometer          int     `json:"startOdometer"`
			StartTime              string  `json:"startTime"`
			StartPosition          struct {
				Longitude       float64 `json:"longitude"`
				Latitude        float64 `json:"latitude"`
				StreetAddress   string  `json:"streetAddress"`
				PostalCode      string  `json:"postalCode"`
				City            string  `json:"city"`
				ISO2CountryCode string  `json:"ISO2CountryCode"`
				Region          string  `json:"Region"`
			} `json:"startPosition"`
			EndOdometer int    `json:"endOdometer"`
			EndTime     string `json:"endTime"`
			EndPosition struct {
				Longitude       float64 `json:"longitude"`
				Latitude        float64 `json:"latitude"`
				StreetAddress   string  `json:"streetAddress"`
				PostalCode      string  `json:"postalCode"`
				City            string  `json:"city"`
				ISO2CountryCode string  `json:"ISO2CountryCode"`
				Region          string  `json:"Region"`
			} `json:"endPosition"`
		} `json:"tripDetails"`
	} `json:"trips"`
	client *Client // added for interface simplification
}

type VehicleServiceStatus struct {
	Status            string  `json:"status"`
	StatusTimestamp   string  `json:"statusTimestamp"`
	StartTime         string  `json:"startTime"`
	ServiceType       string  `json:"serviceType"`
	FailureReason     string  `json:"failureReason"` // TODO: no idea about the actual type
	Service           string  `json:"service"`       // hyperlink
	VehicleID         string  `json:"vehicleId"`     // VIN
	CustomerServiceID string  `json:"customerServiceId"`
	client            *Client // added for interface simplification
}

/*
	The call towards vss.client.Vehicles.GetServiceStatus failes with runtime error: invalid memory address or nil pointer dereference
	which is pretty weird because i'm passing a regular string to a simple function.
*/

func (vss *VehicleServiceStatus) Refresh() error {
	vssNew, err := vss.client.Vehicles.GetServiceStatus(vss.Service)
	if err != nil {
		return err
	}
	// manually refreshing the original struct with the new values
	vss.Status = vssNew.Status
	vss.StatusTimestamp = vssNew.StatusTimestamp
	vss.StartTime = vssNew.StartTime
	vss.ServiceType = vssNew.ServiceType
	vss.FailureReason = vssNew.FailureReason
	vss.Service = vssNew.Service
	vss.VehicleID = vssNew.VehicleID
	vss.CustomerServiceID = vssNew.CustomerServiceID
	return nil
}

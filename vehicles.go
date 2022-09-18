package vocdriver

import "log"

type VehiclesService struct {
	client   *Client
	Endpoint string
}

func (v *VehiclesService) GetVehicleByVIN(vin string) (vehicle *Vehicle, err error) {
	url := v.client.MakeURL(v.Endpoint, vin)
	if _, err = v.client.Request.Get(url, &vehicle); err != nil {
		return nil, err
	}
	vehicle.client = v.client
	return
}

func (v *VehiclesService) GetVehicleByHyperlink(url string) (vehicle *Vehicle, err error) {
	if _, err = v.client.Request.Get(url, &vehicle); err != nil {
		return nil, err
	}
	vehicle.client = v.client
	return
}

func (v *VehiclesService) GetVehicleAttributesByVIN(vin string) (attributes *VehicleAttributes, err error) {
	url := v.client.MakeURL(v.Endpoint, vin, "attributes")
	if _, err = v.client.Request.Get(url, &attributes); err != nil {
		return nil, err
	}
	attributes.client = v.client
	return
}

func (v *VehiclesService) GetVehicleStatusByVIN(vin string) (status *VehicleStatus, err error) {
	url := v.client.MakeURL(v.Endpoint, vin, "status")
	if _, err = v.client.Request.Get(url, &status); err != nil {
		return nil, err
	}
	status.client = v.client
	return
}

func (v *VehiclesService) GetVehiclePositionByVIN(vin string) (position *VehiclePosition, err error) {
	url := v.client.MakeURL(v.Endpoint, vin, "position")
	if _, err = v.client.Request.Get(url, &position); err != nil {
		return nil, err
	}
	position.client = v.client
	return
}

func (v *VehiclesService) GetVehicleTripsByVIN(vin string) (trips *VehicleTrips, err error) {
	url := v.client.MakeURL(v.Endpoint, vin, "trips")
	if _, err = v.client.Request.Get(url, &trips); err != nil {
		return nil, err
	}
	trips.client = v.client
	return
}

func (v *VehiclesService) LockVehicle(vin string) (status *VehicleServiceStatus, err error) {
	url := v.client.MakeURL(v.Endpoint, vin, "lock")
	if _, err = v.client.Request.Post(url, nil, &status); err != nil {
		return nil, err
	}
	return
}

// Utils

func (v *VehiclesService) RetrieveServiceStatus(vin, customerServiceId string) (status *VehicleServiceStatus, err error) {
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
	VehicleAccountRelations          *AccountVehicleRelation
	HyperlinkAttributes              string   `json:"attributes"`              // url
	HyperlinkStatus                  string   `json:"status"`                  // url
	HyperlinkVehicleAccountRelations []string `json:"vehicleAccountRelations"` // url
	VehicleID                        string   `json:"vehicleId"`               // vin
	client                           *Client  // added for interface simplification
}

func (v *Vehicle) RetrieveHyperlinks() (err error) {
	if v.Attributes == nil {
	}
	if v.Status == nil {
	}
	if v.VehicleAccountRelations == nil {
	}
	return
}

func (v *Vehicle) GetAttributes() (attributes *VehicleAttributes, err error) {
	return v.client.Vehicles.GetVehicleAttributesByVIN(v.VehicleID)
}

func (v *Vehicle) GetStatus() (status *VehicleStatus, err error) {
	return v.client.Vehicles.GetVehicleStatusByVIN(v.VehicleID)
}

func (v *Vehicle) GetPosition() (position *VehiclePosition, err error) {
	return v.client.Vehicles.GetVehiclePositionByVIN(v.VehicleID)
}

func (v *Vehicle) GetTrips() (trips *VehicleTrips, err error) {
	return v.client.Vehicles.GetVehicleTripsByVIN(v.VehicleID)
}

func (v *Vehicle) Lock() (status *VehicleServiceStatus, err error) {
	return v.client.Vehicles.LockVehicle(v.VehicleID)
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
	Position struct {
		Longitude float64     `json:"longitude"`
		Latitude  float64     `json:"latitude"`
		Timestamp string      `json:"timestamp"`
		Speed     interface{} `json:"speed"`
		Heading   interface{} `json:"heading"`
	} `json:"position"`
	CalculatedPosition struct {
		Longitude interface{} `json:"longitude"`
		Latitude  interface{} `json:"latitude"`
		Timestamp string      `json:"timestamp"`
		Speed     interface{} `json:"speed"`
		Heading   interface{} `json:"heading"`
	} `json:"calculatedPosition"`
	client *Client // added for interface simplification
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
			FuelConsumption        int    `json:"fuelConsumption"`
			ElectricalConsumption  int    `json:"electricalConsumption"`
			ElectricalRegeneration int    `json:"electricalRegeneration"`
			Distance               int    `json:"distance"`
			StartOdometer          int    `json:"startOdometer"`
			StartTime              string `json:"startTime"`
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
	Service           string  `json:"service"`
	VehicleID         string  `json:"vehicleId"`
	CustomerServiceID string  `json:"customerServiceId"`
	client            *Client // added for interface simplification
}

package vocdriver

type ChargingLocationPosition struct {
	Longitude       float64 `json:"longitude,omitempty"`
	Latitude        float64 `json:"latitude,omitempty"`
	StreetAddress   string  `json:"streetAddress,omitempty"`
	PostalCode      string  `json:"postalCode,omitempty"`
	City            string  `json:"city,omitempty"`
	ISO2CountryCode string  `json:"ISO2CountryCode,omitempty"`
	Region          string  `json:"Region,omitempty"`
}

type DelayCharging struct {
	Enabled   bool   `json:"enabled,omitempty"`
	StartTime string `json:"startTime,omitempty"` // example value: 21:30
	StopTime  string `json:"stopTime,omitempty"`  // example value: 06:45
}

type ChargingLocation struct {
	ChargeLocation            string                    `json:"chargeLocation,omitempty"` // Self URL
	Name                      string                    `json:"name,omitempty"`
	PlugInReminderEnabled     bool                      `json:"plugInReminderEnabled,omitempty"`
	Position                  *ChargingLocationPosition `json:"position,omitempty"`
	DelayCharging             *DelayCharging            `json:"delayCharging,omitempty"`
	Status                    string                    `json:"status,omitempty"`
	VehicleAtChargingLocation bool                      `json:"vehicleAtChargingLocation,omitempty"`
	client                    *Client                   // added for interface simplification
}

type ChargingLocations struct {
	ChargingLocations []ChargingLocation `json:"chargingLocations,omitempty"`
	client            *Client            // added for interface simplification
}

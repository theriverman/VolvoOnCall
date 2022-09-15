package vocdriver

type VehicleAccountRelation struct {
	Account                   string `json:"account"` // url
	AccountID                 string `json:"accountId"`
	Vehicle                   string `json:"vehicle"`                // url
	AccountVehicleRelation    string `json:"accountVehicleRelation"` // url (self)
	VehicleID                 string `json:"vehicleId"`
	Username                  string `json:"username"`
	Status                    string `json:"status"`
	CustomerVehicleRelationID int    `json:"customerVehicleRelationId"`
}

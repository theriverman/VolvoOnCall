package vocdriver

import "fmt"

type VehicleAccountRelationService struct {
	client   *Client
	Endpoint string
}

func (s *VehicleAccountRelationService) GetById(customerVehicleRelationId int) (vehicleAccRel *VehicleAccountRelation, err error) {
	url := s.client.MakeURL(s.Endpoint, fmt.Sprintf("%d", customerVehicleRelationId))
	if _, err = s.client.Request.Get(url, &vehicleAccRel); err != nil {
		return
	}
	return
}

type VehicleAccountRelation struct {
	Account                   string `json:"account"` // url
	AccountID                 string `json:"accountId"`
	Vehicle                   string `json:"vehicle"`                // url
	AccountVehicleRelation    string `json:"accountVehicleRelation"` // url (self)
	VehicleID                 string `json:"vehicleId"`              // VIN
	Username                  string `json:"username"`
	Status                    string `json:"status"`
	CustomerVehicleRelationID int    `json:"customerVehicleRelationId"`
}

package vocdriver

import "fmt"

type AccountVehicleRelationsService struct {
	client   *Client
	Endpoint string
}

func (s *AccountVehicleRelationsService) GetById(customerVehicleRelationId int) (vehicleAccRel *AccountVehicleRelation, err error) {
	url := s.client.MakeURL(s.Endpoint, fmt.Sprintf("%d", customerVehicleRelationId))
	if _, err = s.client.Request.Get(url, &vehicleAccRel); err != nil {
		return
	}
	vehicleAccRel.client = s.client
	return
}

func (s *AccountVehicleRelationsService) GetByHyperlink(url string) (vehicleAccRel *AccountVehicleRelation, err error) {
	if _, err = s.client.Request.Get(url, &vehicleAccRel); err != nil {
		return
	}
	vehicleAccRel.client = s.client
	return
}

type AccountVehicleRelation struct {
	Account                         *CustomerAccount
	Vehicle                         *Vehicle
	AccountVehicleRelation          *AccountVehicleRelation
	VehicleID                       string `json:"vehicleId"`                 // VIN
	Username                        string `json:"username"`                  // typically a phone number
	Status                          string `json:"status"`                    // other states that "Verified" are yet unknown
	CustomerVehicleRelationID       int    `json:"customerVehicleRelationId"` // self primary key
	AccountID                       string `json:"accountId"`                 // uuid
	HyperlinkAccount                string `json:"account"`                   // url
	HyperlinkAccountVehicleRelation string `json:"accountVehicleRelation"`    // url (self)
	HyperlinkVehicle                string `json:"vehicle"`                   // url
	client                          *Client
}

func (avr *AccountVehicleRelation) RetrieveHyperlinks() (err error) {
	if avr.Account == nil {
		if avr.Account, err = avr.client.CustomerAccount.GetAccountByHyperlink(avr.HyperlinkAccount); err != nil {
			return
		}
	}
	if avr.Vehicle == nil {
		if avr.Vehicle, err = avr.client.Vehicles.GetVehicleByHyperlink(avr.HyperlinkVehicle); err != nil {
			return
		}
	}
	if avr.AccountVehicleRelation == nil {
		if avr.AccountVehicleRelation, err = avr.client.AccountVehicleRelation.GetByHyperlink(avr.HyperlinkAccountVehicleRelation); err != nil {
			return
		}
	}
	return
}

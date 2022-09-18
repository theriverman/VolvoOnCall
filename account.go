package vocdriver

import (
	"fmt"
	"strconv"
	"strings"
)

type CustomerAccountService struct {
	client   *Client
	Endpoint string
}

func (s *CustomerAccountService) GetAccount() (customerAccount *CustomerAccount, err error) {
	url := s.client.MakeURL(s.Endpoint)
	if _, err = s.client.Request.Get(url, &customerAccount); err != nil {
		return
	}
	customerAccount.client = s.client
	return
}

func (s *CustomerAccountService) GetAccountByHyperlink(url string) (customerAccount *CustomerAccount, err error) {
	if _, err = s.client.Request.Get(url, &customerAccount); err != nil {
		return
	}
	customerAccount.client = s.client
	return
}

// CustomerAccount is returned at /customeraccounts
type CustomerAccount struct {
	AccountVehicleRelations          []AccountVehicleRelation
	Username                         string   `json:"username"` // (can be phone number)
	FirstName                        string   `json:"firstName"`
	LastName                         string   `json:"lastName"`
	AccountID                        string   `json:"accountId"`               // uuid
	HyperlinkAccount                 string   `json:"account"`                 // url
	HyperlinkAccountVehicleRelations []string `json:"accountVehicleRelations"` // urls
	client                           *Client  // added for interface simplification
	accountVehicleRelationsRetrieved bool
}

func (ca *CustomerAccount) RetrieveHyperlinks() (err error) {
	if !ca.accountVehicleRelationsRetrieved {
		for _, url := range ca.HyperlinkAccountVehicleRelations {
			relation, err := ca.client.AccountVehicleRelation.GetByHyperlink(url)
			if err != nil {
				return err
			}
			ca.AccountVehicleRelations = append(ca.AccountVehicleRelations, *relation)
		}
		ca.accountVehicleRelationsRetrieved = true
	}
	return
}

func (ca CustomerAccount) GetAccountVehicleRelationIds() (relationIds []int, err error) {
	for _, url := range ca.HyperlinkAccountVehicleRelations {
		urlParts := strings.Split(url, "/")
		relationId, err := strconv.Atoi(urlParts[len(urlParts)-1])
		if err != nil {
			return relationIds, fmt.Errorf("failed to extract relation id from %s", url)
		}
		relationIds = append(relationIds, relationId)
	}
	return
}

func (ca *CustomerAccount) GetAccountVehicleRelations() (relations []AccountVehicleRelation, err error) {
	relationIds, err := ca.GetAccountVehicleRelationIds()
	if err != nil {
		return
	}
	for _, id := range relationIds {
		relation, err := ca.client.AccountVehicleRelation.GetById(id)
		if err != nil {
			return relations, err
		}
		relations = append(relations, *relation)
	}
	return
}

func (ca *CustomerAccount) GetVehicles() (vehicles []Vehicle, err error) {
	accountVehicleRelations, err := ca.GetAccountVehicleRelations()
	if err != nil {
		return
	}

	for _, relation := range accountVehicleRelations {
		vehicle, err := ca.client.Vehicles.GetVehicleByVIN(relation.VehicleID)
		if err != nil {
			return vehicles, err
		}
		vehicles = append(vehicles, *vehicle)
	}
	return
}

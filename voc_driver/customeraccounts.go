package vocdriver

import (
	"log"
	"strconv"
	"strings"
)

type CustomerAccountsService struct {
	client   *Client
	Endpoint string
}

func (s *CustomerAccountsService) GetAccount() (customerAccount *CustomerAccount, err error) {
	url := s.client.MakeURL(s.Endpoint)
	if _, err = s.client.Request.Get(url, &customerAccount); err != nil {
		return
	}
	return
}

// CustomerAccount is returned at /customeraccounts
type CustomerAccount struct {
	Username                string   `json:"username"` // (can be phone number)
	FirstName               string   `json:"firstName"`
	LastName                string   `json:"lastName"`
	AccountID               string   `json:"accountId"`               // uuid
	Account                 string   `json:"account"`                 // url
	AccountVehicleRelations []string `json:"accountVehicleRelations"` // urls
}

func (ca CustomerAccount) GetAccountVehicleRelations() (accountVehicleRelationsIds []int) {
	for _, e := range ca.AccountVehicleRelations {
		urlParts := strings.Split(e, "/")
		customerVehicleRelationId, err := strconv.Atoi(urlParts[len(urlParts)-1])
		if err != nil {
			log.Panicln(err)
		}
		accountVehicleRelationsIds = append(accountVehicleRelationsIds, customerVehicleRelationId)
	}
	return
}

package vocdriver

// CustomerAccount is returned at /customeraccounts
type CustomerAccount struct {
	Username                string   `json:"username"`
	FirstName               string   `json:"firstName"`
	LastName                string   `json:"lastName"`
	AccountID               string   `json:"accountId"`
	Account                 string   `json:"account"`                 // url (self)
	AccountVehicleRelations []string `json:"accountVehicleRelations"` // url
}

type CustomerAccountsService struct {
	client   *Client
	Endpoint string
}

func (s *CustomerAccountsService) Get() (*CustomerAccount, error) {
	url := s.client.MakeURL(s.Endpoint)
	var ca CustomerAccount
	_, err := s.client.Request.Get(url, &ca)
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

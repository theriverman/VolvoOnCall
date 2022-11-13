package vocdriver

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

/*
Possible ServiceRegion options:
- Africa
- Asian (except China Mainland)
- China Mainland
- Europe [default]
- North America
- Oceania
- South America
*/

const BaseUrl string = "https://vocapi%s.wirelesscar.net/customerapi/rest/v3.0"

func NewClient(username, password string) (*Client, error) {
	client := Client{}
	client.Initialise()
	client.Authenticate(username, password)
	return &client, nil
}

type Client struct {
	apiUrl        string
	BaseURL       string
	ServiceRegion string
	Headers       *http.Header
	HTTPClient    *http.Client

	isInitialised bool
	Verbose       bool

	// Core Services
	Request *RequestService

	// Volvo Services
	CustomerAccount        *CustomerAccountService
	AccountVehicleRelation *AccountVehicleRelationsService
	Vehicles               *VehiclesService
}

func (c *Client) Initialise() error {
	// skip if already Initialised
	if c.isInitialised {
		return nil
	}

	// add a default http.Client{}
	c.HTTPClient = &http.Client{
		// CheckRedirect: redirectPolicyFunc,
	}

	// set basic headers
	c.Headers = &http.Header{}
	c.Headers.Add("Content-Type", "application/json")
	c.Headers.Add("X-App-Name", "Volvo On Call")
	c.Headers.Add("X-Client-Version", "4.4.5.21126")
	c.Headers.Add("X-Device-Id", "Device")
	c.Headers.Add("X-Originator-Type", "App")
	c.Headers.Add("X-OS-Type", "Android")
	c.Headers.Add("X-OS-Version", "22")

	switch {
	case c.BaseURL == "" && c.ServiceRegion == "":
		c.apiUrl = fmt.Sprintf(BaseUrl, c.ServiceRegion)
	case c.BaseURL == "" && c.ServiceRegion != "":
		c.apiUrl = fmt.Sprintf(c.BaseURL, "-"+c.ServiceRegion)
	case c.BaseURL != "":
		c.apiUrl = c.BaseURL // ServiceRegion must be defined as part of the BaseUrl
	}

	// Bootstrapping Services
	c.Request = &RequestService{c}
	c.CustomerAccount = &CustomerAccountService{c, "customeraccounts"}
	c.Vehicles = &VehiclesService{c, "vehicles"}
	c.AccountVehicleRelation = &AccountVehicleRelationsService{c, "vehicle-account-relations"}
	c.isInitialised = true
	return nil
}

// LoadExternalHeaders loads a map of header key/value pairs permemently into `Client.Headers`
func (c *Client) LoadExternalHeaders(headers map[string]string) {
	for k, v := range headers {
		c.Headers.Add(k, v)
	}
}

// Authenticate encodes username+password using base64 and adds the resulting string to default Headers
func (c *Client) Authenticate(username, password string) {
	c.Headers.Add("Authorization", "Basic "+basicAuth(username, password))
}

// MakeURL accepts an Endpoint URL and returns a compiled absolute URL
//
//	For example:
//	- If the given endpoint URLs are [epics, attachments]
//	- If the BaseURL is https://api.taiga.io
//	- It returns https://api.taiga.io/api/v1/epics/attachments
//	- Suffixes are appended to the URL joined by a slash (/)
func (c *Client) MakeURL(EndpointParts ...string) string {
	return c.apiUrl + "/" + strings.Join(EndpointParts, "/")
}

func (c *Client) EvaluateServiceStatus(vss *VehicleServiceStatus, timeoutSeconds int) (err error) {
	return c.Vehicles.EvaluateServiceStatus(vss, timeoutSeconds)
}

func (c *Client) EvaluateServiceStatusAuto(vss *VehicleServiceStatus) (err error) {
	return c.Vehicles.EvaluateServiceStatusAuto(vss)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

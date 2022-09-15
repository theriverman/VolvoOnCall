package vocdriver

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

const BaseUrl string = "https://vocapi%s.wirelesscar.net/customerapi/rest/v3.0"

type Client struct {
	apiUrl        string
	BaseURL       string
	Headers       *http.Header
	HTTPClient    *http.Client
	ServiceRegion string

	isInitialised bool
	Verbose       bool

	// Core Services
	Request *RequestService

	// Volvo Services
	CustomerAccounts *CustomerAccountsService
	// Auth      *AuthService
	// Epic      *EpicService
	// Issue     *IssueService
	// Milestone *MilestoneService
	// Project   *ProjectService
	// Resolver  *ResolverService
	// Stats     *StatsService
	// Task      *TaskService
	// UserStory *UserStoryService
	// User      *UserService
	// Webhook   *WebhookService
	// Wiki      *WikiService

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

	// check baseUrl
	if c.BaseURL == "" {
		c.apiUrl = fmt.Sprintf(BaseUrl, c.ServiceRegion)
	} else {
		c.apiUrl = fmt.Sprintf(c.BaseURL, c.ServiceRegion)
	}

	// Bootstrapping Services
	c.Request = &RequestService{c}
	c.CustomerAccounts = &CustomerAccountsService{c, "customeraccounts"}
	// c.Auth = &AuthService{c, 0, "auth"}
	// c.Epic = &EpicService{c, 0, "epics"}
	// c.Issue = &IssueService{c, 0, "issues"}
	// c.Milestone = &MilestoneService{c, 0, "milestones"}
	// c.Project = &ProjectService{client: c, Endpoint: "projects"}
	// c.Resolver = &ResolverService{c, 0, "resolver"}
	// c.Stats = &StatsService{c, 0, "stats"}
	// c.Task = &TaskService{c, 0, "tasks"}
	// c.UserStory = &UserStoryService{c, 0, "userstories"}
	// c.User = &UserService{c, 0, "users"}
	// c.Webhook = &WebhookService{c, 0, "webhooks", "webhooklogs"}
	// c.Wiki = &WikiService{c, 0, "wiki"}

	c.isInitialised = true
	return nil
}

func (c *Client) loadHeaders(request *http.Request) {
	for key, values := range *c.Headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
}

// func (c *Client) setToken() {
// 	c.Headers.Del("Authorization") // avoid header duplication
// 	c.Headers.Add("Authorization", c.TokenType+" "+c.Token)
// }

// LoadExternalHeaders loads a map of header key/value pairs permemently into `Client.Headers`
func (c *Client) LoadExternalHeaders(headers map[string]string) {
	for k, v := range headers {
		c.Headers.Add(k, v)
	}
}

func (c *Client) Authenticate(username, password string) {
	c.Headers.Add("Authorization", "Basic "+basicAuth(username, password))
}

// MakeURL accepts an Endpoint URL and returns a compiled absolute URL

// For example:
//   - If the given endpoint URLs are [epics, attachments]
//   - If the BaseURL is https://api.taiga.io
//   - It returns https://api.taiga.io/api/v1/epics/attachments
//   - Suffixes are appended to the URL joined by a slash (/)
func (c *Client) MakeURL(EndpointParts ...string) string {
	return c.apiUrl + "/" + strings.Join(EndpointParts, "/")
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
// 	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
// 	return nil
// }

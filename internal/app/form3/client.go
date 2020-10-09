package form3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Params struct {
	Page     *int
	PageSize *int
}

type HttpClient interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	http HttpClient
	Host string
}

// NewClient returns a new Client used to make request to Account API
func NewClient(http HttpClient, h string) *Client {
	return &Client{
		http: http,
		Host: h,
	}
}

// Create creates a new bank account
func (c *Client) Create(ac Account) (*Account, error) {
	body, err := json.Marshal(map[string]Account{"data": ac})
	if err != nil {
		return nil, err
	}

	URL := fmt.Sprintf("%s/v1/organisation/accounts", c.Host)
	res, err := c.http.Post(URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not make request: %s", err))
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read response: %s", err))
	}

	var account struct {
		Data Account `json:"data"`
	}

	err = json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}

	return &account.Data, nil
}

// Fetch returns a bank account with the provided id
func (c *Client) Fetch(id string) (*Account, error) {
	URL := fmt.Sprintf("%s/v1/organisation/accounts/%s", c.Host, id)
	res, err := c.http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read response: %s", err))
	}

	var account struct {
		Data Account `json:"data"`
	}

	err = json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}

	return &account.Data, nil
}

// List returns a list of bank accounts paginated according with the provided parameters
func (c *Client) List(params Params) ([]Account, error) {
	URL := fmt.Sprintf("%s/v1/organisation/accounts", c.Host)
	u, _ := url.Parse(URL)
	q := u.Query()

	if params.Page != nil {
		p := *params.Page
		q.Set("page[number]", strconv.Itoa(p))
	}

	if params.PageSize != nil {
		pz := *params.PageSize
		q.Set("page[size]", strconv.Itoa(pz))
	}

	u.RawQuery = q.Encode()

	res, err := c.http.Get(u.String())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read response: %s", err))
	}

	var accounts struct {
		Data []Account `json:"data"`
	}

	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts.Data, nil
}

// Delete removes the version of a bank account with the provided id
func (c *Client) Delete(id string, v int) error {
	URL := fmt.Sprintf("%s/v1/organisation/accounts/%s", c.Host, id)

	u, _ := url.Parse(URL)
	q := u.Query()
	q.Set("version", strconv.Itoa(v))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create request: %s", err))
	}

	_, err = c.http.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to make request: %s", err))
	}

	return nil
}

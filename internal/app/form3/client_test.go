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
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var host = "http://localhost:8080"

// TestMain main run function
func TestMain(m *testing.M) {
	if val, present := os.LookupEnv("ACCOUNT_API_HOST"); present {
		host = val
	}

	os.Exit(m.Run())
}

func TestClient_Create(t *testing.T) {
	a := getDummyAccounts()[0]

	testCases := []struct {
		name         string
		requestBody  interface{}
		responseBody interface{}
		input        Account
		err          error
	}{
		{
			name:         "account is created successfully",
			requestBody:  map[string]Account{"data": a},
			responseBody: map[string]Account{"data": a},
			input:        a,
			err:          nil,
		},
		{
			name:         "account api request error",
			requestBody:  map[string]Account{"data": a},
			responseBody: nil,
			input:        a,
			err:          errors.New("service is down"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			URL := fmt.Sprintf("%s/v1/organisation/accounts", host)

			reqBody, _ := json.Marshal(tc.requestBody)
			resBody, _ := json.Marshal(tc.responseBody)
			mockRes := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(resBody))}

			mockHttpClient := HttpClientMock{}
			mockHttpClient.On("Post", URL, "application/json", bytes.NewBuffer(reqBody)).Return(mockRes, tc.err)

			form3Client := NewClient(&mockHttpClient, host)
			newAccount, err := form3Client.Create(tc.input)

			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, newAccount)
				assert.True(t, len(newAccount.ID) > 0)
			}
		})
	}
}

func TestClient_Fetch(t *testing.T) {
	a := getDummyAccounts()[0]

	testCases := []struct {
		name         string
		requestURL   string
		responseBody interface{}
		input        string
		err          error
	}{
		{
			name:         "account fetched successfully",
			requestURL:   fmt.Sprintf("%s/v1/organisation/accounts", host),
			responseBody: map[string]Account{"data": a},
			input:        "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			err:          nil,
		},
		{
			name:         "account api request error",
			requestURL:   fmt.Sprintf("%s/v1/organisation/accounts", host),
			responseBody: nil,
			input:        "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			err:          errors.New("service is down"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			resBody, _ := json.Marshal(tc.responseBody)
			mockRes := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(resBody))}

			mockHttpClient := HttpClientMock{}
			URL := tc.requestURL + "/" + tc.input
			mockHttpClient.On("Get", URL).Return(mockRes, tc.err)

			form3Client := NewClient(&mockHttpClient, host)
			account, err := form3Client.Fetch(tc.input)

			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, a.ID, account.ID)
			}
		})
	}
}

func TestClient_List(t *testing.T) {
	a, p, pz := getDummyAccounts()[0], 1, 10

	testCases := []struct {
		name         string
		requestURL   string
		responseBody interface{}
		params       Params
		err          error
	}{
		{
			name:         "accounts fetched successfully without query params",
			requestURL:   fmt.Sprintf("%s/v1/organisation/accounts", host),
			responseBody: map[string][]Account{"data": {a}},
			params:       Params{},
			err:          nil,
		},
		{
			name:         "accounts fetched successfully with query params",
			requestURL:   fmt.Sprintf("%s/v1/organisation/accounts", host),
			responseBody: map[string][]Account{"data": {a}},
			params:       Params{Page: &p, PageSize: &pz},
			err:          nil,
		},
		{
			name:         "account api request error",
			requestURL:   fmt.Sprintf("%s/v1/organisation/accounts", host),
			responseBody: nil,
			err:          errors.New("service is down"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			resBody, _ := json.Marshal(tc.responseBody)
			mockRes := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer(resBody))}

			u, _ := url.Parse(tc.requestURL)
			q := u.Query()
			if tc.params.Page != nil {
				p := *tc.params.Page
				q.Set("page[number]", strconv.Itoa(p))
			}

			if tc.params.PageSize != nil {
				pz := *tc.params.PageSize
				q.Set("page[size]", strconv.Itoa(pz))
			}
			u.RawQuery = q.Encode()

			mockHttpClient := HttpClientMock{}
			mockHttpClient.On("Get", u.String()).Return(mockRes, tc.err)

			form3Client := NewClient(&mockHttpClient, host)
			accounts, err := form3Client.List(tc.params)

			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, len(accounts) == 1)
			}
		})
	}
}

func TestClient_Delete(t *testing.T) {
	testCases := []struct {
		name       string
		requestURL string
		input      string
		version    int
		err        error
	}{
		{
			name:       "account deleted successfully",
			requestURL: fmt.Sprintf("%s/v1/organisation/accounts", host),
			input:      "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:    0,
			err:        nil,
		},
		{
			name:       "account api request error",
			requestURL: fmt.Sprintf("%s/v1/organisation/accounts", host),
			input:      "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			version:    0,
			err:        errors.New("service is down"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			mockHttpClient := HttpClientMock{}
			URL := tc.requestURL + "/" + tc.input + "?version=" + strconv.Itoa(tc.version)
			req, err := http.NewRequest(http.MethodDelete, URL, nil)
			mockResponse := &http.Response{Body: nil}
			mockHttpClient.On("Do", req).Return(mockResponse, tc.err)

			form3Client := NewClient(&mockHttpClient, host)
			err = form3Client.Delete(tc.input, tc.version)

			if tc.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func getDummyAccounts() []Account {
	return []Account{
		{
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			Type:           "accounts",
			Attributes: Attributes{
				Country:      "GB",
				BaseCurrency: "GBP",
				BankID:       "400300",
				BankIDCode:   "GBDSC",
				BIC:          "NWBKGB22",
			},
		},
		{
			ID:             "da171a64-f1db-40ff-a29f-f852ce95d7b1",
			OrganisationID: "6224e61c-f135-4708-9b4d-3afefa6f9c76",
			Type:           "accounts",
			Attributes: Attributes{
				Country:      "GR",
				BaseCurrency: "EUR",
				BankID:       "400300",
				BankIDCode:   "GRBIC",
				BIC:          "NWBKGB22",
			},
		},
		{
			ID:             "e2a1b488-9d11-4d5c-8a6e-1a4918c6ffbf",
			OrganisationID: "206133da-d3dd-4665-8a88-fc8e720dae3e",
			Type:           "accounts",
			Attributes: Attributes{
				Country:      "FR",
				BaseCurrency: "EUR",
				BankID:       "400300",
				BankIDCode:   "CAPCA",
				BIC:          "NWBKGB22",
			},
		},
	}
}

type HttpClientMock struct {
	mock.Mock
}

func (c *HttpClientMock) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	args := c.Called(url, contentType, body)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *HttpClientMock) Get(url string) (*http.Response, error) {
	args := c.Called(url)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}

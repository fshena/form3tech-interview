// +build integration

package form3

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateAccount(t *testing.T) {
	client := NewClient(&http.Client{}, host)
	accounts := getDummyAccounts()

	for i := 0; i < len(accounts); i++ {
		a, err := client.Create(accounts[i])
		assert.NoError(t, err)
		assert.True(t, len(a.ID) > 0)
	}
}

func TestClient_GetAccount(t *testing.T) {
	client := NewClient(&http.Client{}, host)
	a, err := client.Fetch(getDummyAccounts()[0].ID)

	assert.NoError(t, err)
	assert.True(t, len(a.ID) > 0)
}

func TestClient_GetAccountWithPagination(t *testing.T) {
	client := NewClient(&http.Client{}, host)
	p, ps := 1, 1
	a, err := client.List(Params{Page: &p, PageSize: &ps})

	assert.NoError(t, err)
	assert.True(t, len(a) == 1)
	assert.Equal(t, getDummyAccounts()[1].ID, a[0].ID)
}

func TestClient_DeleteAccount(t *testing.T) {
	accounts := getDummyAccounts()
	client := NewClient(&http.Client{}, host)

	for i := 0; i < len(accounts); i++ {
		err := client.Delete(accounts[i].ID, 0)
		assert.NoError(t, err)
	}
}

package bunq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const monetaryAccountID int = 9601

func TestMonetaryAccountBankListing(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := c.AccountService.GetAllMonetaryAccountBank()

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].MonetaryAccountBank.ID)
}

func TestMonetaryAccountBankGet(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := c.AccountService.GetMonetaryAccountBank(monetaryAccountID)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].MonetaryAccountBank.ID)
}

func TestAccountService_GetAllMonetaryAccountSaving(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := c.AccountService.GetAllMonetaryAccountSaving()

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].MonetaryAccountSaving.ID)
}

func TestAccountService_GetMonetaryAccountSaving(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := c.AccountService.GetMonetaryAccountSaving(monetaryAccountID)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].MonetaryAccountSaving.ID)
}

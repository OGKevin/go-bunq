package bunq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduleListing(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := c.ScheduledPaymentService.GetAllScheduledPayments(monetaryAccountID)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].ScheduledPayment.MonetaryAccountID)
}

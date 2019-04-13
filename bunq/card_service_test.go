package bunq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCardService_GetMasterCardAction(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := c.cardService.GetMasterCardAction(324, 9520)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].MasterCardAction.ID)
}

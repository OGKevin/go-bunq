package bunq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceServerCreate(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	_, err := c.installation.create()
	assert.NoError(t, err)
	res, err := c.deviceServer.create()

	assert.NoError(t, err)
	assert.Equal(t, getDeviceServerResponse(t), res)
}

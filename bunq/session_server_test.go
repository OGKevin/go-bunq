package bunq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionServerCreate(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()
	_, err := c.installation.create()
	assert.NoError(t, err)

	_, err = c.deviceServer.create()
	assert.NoError(t, err)

	r, err := c.sessionServer.create()
	assert.NoError(t, err)

	assert.Equal(t, r.Response[0].Token.Token, *c.token)
}

func TestSessionDeleteOnCancel(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	cancel()

	time.Sleep(time.Second * 2)
}

package bunq

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	fakeServer := httptest.NewServer(createBunqFakeHandlerWithError(t, "/v1/device-server"))
	defer fakeServer.Close()

	key, err := CreateNewKeyPair()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := NewClient(ctx, fmt.Sprintf("%s/v1/", fakeServer.URL), key, "")

	_, err = c.installation.create()
	assert.NoError(t, err)

	_, err = c.deviceServer.create()
	assert.Error(t, err)
}

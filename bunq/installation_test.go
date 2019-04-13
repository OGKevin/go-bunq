package bunq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntsallation(t *testing.T) {
	t.Parallel()

	resInstallation := getInstallationResponse(t)

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	installationRespActual, err := c.installation.create()
	if !assert.NoError(t, err ) {
		return
	}

	installationRespProp := responseInstallation{
		Response: []installation{
			{
				ID:              resInstallation.Response[0].ID,
				Token:           resInstallation.Response[1].Token,
				ServerPublicKey: resInstallation.Response[2].ServerPublicKey,
			},
		},
	}

	assert.EqualValues(t, &installationRespProp, installationRespActual)
	assert.Equal(t, installationRespProp.Response[0].Token.Token, *c.token, nil)
}

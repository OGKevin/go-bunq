package bunq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClientContextExportAndImport(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	exportedCtx, err := c.ExportClientContext()

	assert.NoError(t, err)

	ctx, cancl := context.WithCancel(context.Background())
	defer cancl()
	cFromCtx, err := NewClientFromContext(ctx, &exportedCtx)

	assert.NoError(t, err)
	assert.Equal(t, c.apiKey, cFromCtx.apiKey)
	assert.Equal(t, c.baseURL, cFromCtx.baseURL)
	assert.Equal(t, c.privateKey, cFromCtx.privateKey)
	assert.Equal(t, c.serverPublicKey, cFromCtx.serverPublicKey)
	assert.Equal(t, c.token, cFromCtx.token)
	assert.Equal(t, c.installationContext, cFromCtx.installationContext)
	assert.Equal(t, c.sessionServerContext, cFromCtx.sessionServerContext)
	assert.True(t, c.IsUserPerson())
	assert.True(t, cFromCtx.IsUserPerson())
	assert.False(t, c.IsUserCompany())
	assert.False(t, c.IsUserAPIKey())
	assert.False(t, cFromCtx.IsUserCompany())
	assert.False(t, cFromCtx.IsUserAPIKey())
}

func TestSessionRenewal(t *testing.T) {
	t.SkipNow()
	_, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	var cContext ClientContext

	file, err := os.Open("../testdata/bunq/client_context.json")
	if err != nil {
		t.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&cContext)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancl := context.WithCancel(context.Background())
	defer cancl()

	c, err := NewClientFromContext(ctx, &cContext)
	if err != nil {
		t.Fatal(err)
	}

	c.baseURL = fmt.Sprintf("%s/v1/", fakeServer.URL)
	c.sessionServerContext.UserPerson.SessionTimeout = 1

	sessionToken := c.sessionServerContext.Token.Token
	c.tokenMutex.RLock()
	token := *c.token
	c.tokenMutex.RUnlock()

	assert.NoError(t, c.Init())

	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()

	assert.NotEqual(t, sessionToken, c.sessionServerContext.Token.Token)
	assert.NotEqual(t, token, c.token)
	assert.Equal(t, *c.token, c.sessionServerContext.Token.Token)
}

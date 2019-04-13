package bunq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_contentService_GetAttachmentPublic(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	s, err := c.ContentService.GetAttachmentPublic("f9a1a89a-fdc1-4de5-89d5-e477cccd22c4")
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

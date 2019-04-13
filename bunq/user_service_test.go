package bunq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserPerson(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	r, err := c.UserService.GetUserPerson()

	assert.NoError(t, err)
	assert.NotZero(t, r.Response[0].UserPerson.ID)
}

func TestUpdateUserPerson(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	bod := requestUserPersonPut{
		NotificationFilters: []notificationFilter{
			{
				NotificationDeliveryMethod: "URL",
				NotificationTarget:         "https://requestbin.fullcontact.com/pwgm46pw",
				Category:                   "MUTATION",
			},
		},
	}

	res, err := c.UserService.UpdateUserPerson(bod)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].ID.ID)
}

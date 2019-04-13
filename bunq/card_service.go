package bunq

import (
	"fmt"
	"net/http"
)

type cardService service

func (c *cardService) GetMasterCardAction(id, monetaryAccountID int) (*responseMasterCardActionGet, error) {
	userID, err := c.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := c.client.preformRequest(http.MethodGet, c.client.formatRequestURL(fmt.Sprintf(endpointMasterCardActionGet, userID, monetaryAccountID, id)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct responseMasterCardActionGet

	return &resStruct, c.client.parseResponse(res, &resStruct)
}

package bunq

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type schedulePaymentService service

func (sp *schedulePaymentService) GetAllSchedulePayment(monetaryAccountID int) (*ResponseSchedulePaymentGet, error) {
	userID, err := sp.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := sp.client.preformRequest(http.MethodGet, sp.client.formatRequestURL(fmt.Sprintf(endpointSchedulePaymentGet, userID, monetaryAccountID)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to get all scheduled payments failed")
	}

	var resSpGet ResponseSchedulePaymentGet

	return &resSpGet, sp.client.parseResponse(res, &resSpGet)
}

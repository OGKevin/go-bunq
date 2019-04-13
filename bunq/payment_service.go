package bunq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type paymentService service

func (p *paymentService) CreateDraftPayment(monetaryAccountID int, rBody requestCreateDraftPayment) (*responseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(rBody)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointDraftPaymentCreate, userID, monetaryAccountID)), bodyRaw, http.MethodPost)
}

func (p *paymentService) UpdateDraftPayment(id, monetaryAccountID int, rBody requestUpdateDraftPayment) (*responseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(rBody)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointDraftPaymentWithID, userID, monetaryAccountID, id)), bodyRaw, http.MethodPut)
}

func (p *paymentService) GetDraftPayment(id, monetaryAccountID int) (*responseDraftPaymentGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointDraftPaymentWithID, userID, monetaryAccountID, id)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct responseDraftPaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetPayment returns a specific payment for a given account
func (p *paymentService) GetPayment(monetaryAccountID uint, paymentID uint) (*ResponsePaymentGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, errors.Wrap(err, "bunq: payment service: could not determine user id")
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointPaymentGetWithID, userID, monetaryAccountID, paymentID)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct ResponsePaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetAllPayment returns all the payments for a given account
func (p *paymentService) GetAllPayment(monetaryAccountID uint) (*ResponsePaymentGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, errors.Wrap(err, "bunq: payment service: could not determine user id")
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointPaymentGet, userID, monetaryAccountID)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct ResponsePaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetAllOlderPayment calls the older url from the Pagination
func (p *paymentService) GetAllOlderPayment(pagi Pagination) (*ResponsePaymentGet, error) {
	if pagi.OlderURL == "" {
		return nil, nil
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(pagi.OlderURL[len("/v1/"):]), nil)
	if err != nil {
		return nil, err
	}

	var resStruct ResponsePaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

func (p *paymentService) createPaymentBatch(monetaryAccountID int, create PaymentBatchCreate) (*responseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(create)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointPaymentBatchCreate, userID, monetaryAccountID)), bodyRaw, http.MethodPost)
}

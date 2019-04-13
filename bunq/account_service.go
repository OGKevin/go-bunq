package bunq

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type accountService service

func (a *accountService) GetAllMonetaryAccountBank() (*ResponseMonetaryAccountBankGet, error) {
	userID, err := a.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := a.client.preformRequest(http.MethodGet, a.client.formatRequestURL(fmt.Sprintf(endpointMonetaryAccountBankListing, userID)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to get all MA bank failed")
	}

	var resMaGet ResponseMonetaryAccountBankGet

	return &resMaGet, a.client.parseResponse(res, &resMaGet)
}

func (a *accountService) GetMonetaryAccountBank(id int) (*ResponseMonetaryAccountBankGet, error) {
	userID, err := a.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := a.client.preformRequest(http.MethodGet, a.client.formatRequestURL(fmt.Sprintf(endpointMonetaryAccountBankGet, userID, id)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to get MA bank failed")
	}

	var resMaGet ResponseMonetaryAccountBankGet

	return &resMaGet, a.client.parseResponse(res, &resMaGet)
}

func (a *accountService) GetAllMonetaryAccountSaving() (*ResponseMonetaryAccountSavingGet, error) {
	userID, err := a.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := a.client.preformRequest(http.MethodGet, a.client.formatRequestURL(fmt.Sprintf(endpointMonetaryAccountSavingsListing, userID)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to get all MA saving failed")
	}

	var resStruct ResponseMonetaryAccountSavingGet

	return &resStruct, a.client.parseResponse(res, &resStruct)
}

func (a *accountService) GetMonetaryAccountSaving(id int) (*ResponseMonetaryAccountSavingGet, error) {
	userID, err := a.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := a.client.preformRequest(http.MethodGet, a.client.formatRequestURL(fmt.Sprintf(endpointMonetaryAccountSavingsGet, userID, id)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to get MA saving failed")
	}

	var resStruct ResponseMonetaryAccountSavingGet

	return &resStruct, a.client.parseResponse(res, &resStruct)
}

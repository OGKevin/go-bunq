package bunq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type userService service

// GetUserPerson retrieves a signle user person. Because there can be 1 user person per api key.
// the user id will be determined by the client.
// https://doc.bunq.com/#/user-person/Read_UserPerson
func (u *userService) GetUserPerson() (*responseUserPerson, error) {
	userID, err := u.client.GetUserID()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodGet, u.client.formatRequestURL(fmt.Sprintf(endpointUserPersonGet, userID)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not create request for user-person")
	}

	res, err := u.client.do(r)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to user-person failed")
	}

	var resUserPerson responseUserPerson
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resUserPerson)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: parsing response failed")
	}

	return &resUserPerson, nil
}

// UpdateUserPerson updates the contents of the current auth user-person.
// https://doc.bunq.com/#/user-person/Update_UserPerson
func (u *userService) UpdateUserPerson(rBody requestUserPersonPut) (*responseBunqID, error) {
	userID, err := u.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(rBody)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	r, err := http.NewRequest(
		http.MethodPut,
		u.client.formatRequestURL(fmt.Sprintf(endpointUserPersonGet, userID)),
		bytes.NewBuffer(bodyRaw),
	)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not create request for user-person")
	}

	res, err := u.client.do(r)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to user-person failed")
	}

	var resBunqID responseBunqID
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resBunqID)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: parsing response failed")
	}

	return &resBunqID, nil
}

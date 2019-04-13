package bunq

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
	"net/http"
)

type installationService service

func (i installationService) create() (*responseInstallation, error) {
	body, err := i.createInstallationBody()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(
		http.MethodPost,
		i.client.formatRequestURL(endpointInstallationCreate),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	res, err := i.client.do(r)
	if err != nil {
		return nil, err
	}

	var resInstallation responseInstallation
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resInstallation)

	properInstalltionResponse := createProperInstallationResponse(resInstallation)
	i.setInstallationContextToClient(properInstalltionResponse)

	return &properInstalltionResponse, err
}

func (i *installationService) createInstallationBody() ([]byte, error) {
	if i.client.privateKey == nil {
		return nil, errors.New("bunq: private key has not been set")
	}

	pubKey, err := x509.MarshalPKIXPublicKey(i.client.privateKey.Public())
	if err != nil {
		return nil, err
	}

	pubKeyPemBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubKey,
	}

	pubBtye := pem.EncodeToMemory(&pubKeyPemBlock)
	bodyStruct := requestInstallation{ClientPublicKey: string(pubBtye)}
	bodyData, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, err
	}

	return bodyData, nil
}

func createProperInstallationResponse(res responseInstallation) responseInstallation {
	return responseInstallation{
		Response: []installation{
			{
				ID:              res.Response[0].ID,
				Token:           res.Response[1].Token,
				ServerPublicKey: res.Response[2].ServerPublicKey,
			},
		},
	}
}

func (i *installationService) setInstallationContextToClient(res responseInstallation) {
	i.client.installationContext = &res.Response[0]

	i.client.tokenMutex.Lock()
	defer i.client.tokenMutex.Unlock()
	i.client.token = &i.client.installationContext.Token.Token

	block, _ := pem.Decode([]byte(res.Response[0].ServerPublicKey.ServerPublicKey))
	parseResult, _ := x509.ParsePKIXPublicKey(block.Bytes)
	key := parseResult.(*rsa.PublicKey)
	i.client.serverPublicKey = key
}

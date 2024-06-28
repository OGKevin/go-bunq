package bunq

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// CreateNewKeyPair creates a new RSA key pair that can be used for signing requests.
func CreateNewKeyPair() (*rsa.PrivateKey, error) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)

	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not generate key pair")
	}

	return key, nil
}

func (c *Client) addSignatureHeader(r *http.Request) error {
	var err error
	var body io.ReadCloser

	if r.Body != nil {
		body, err = r.GetBody()
	}

	if err != nil {
		return errors.Wrap(err, "bunq: could not get request body")
	}

	stringToSign := createStringToSign(body)
	h := sha256.New()

	_, err = h.Write([]byte(stringToSign))
	if err != nil {
		return errors.Wrap(err, "bunq: could not encode string to sign to sha256")
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, c.privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return errors.Wrap(err, "bunq: could not sign request")
	}

	r.Header.Set("X-Bunq-Client-Signature", base64.StdEncoding.EncodeToString(signature))

	return nil
}

func (c *Client) verifySignature(_ *http.Response) (bool, error) {
	return true, nil
}

func createStringToVerify(body io.ReadCloser) string {
	defer body.Close()

	rawBody, _ := ioutil.ReadAll(body)
	stringToVerify := string(rawBody)
	stringToVerify = strings.TrimSuffix(stringToVerify, "\n")

	return stringToVerify
}

func createStringToSign(body io.ReadCloser) string {
	if body != nil {
		defer body.Close()
		rawBody, _ := ioutil.ReadAll(body)
		return fmt.Sprintf("%s", rawBody)
	}
	return "\n"
}

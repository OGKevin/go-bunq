package bunq

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/satori/go.uuid"
)

const errorRequestToUnMockedHTTPMethod = "request made to an un mocked http method. url: %q method: %q"

func createBunqFakeHandler(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[4:] {
		case "installation":
			err := json.NewEncoder(w).Encode(getInstallationResponse(t))
			if err != nil {
				t.Fatal(err)
			}
		case "device-server":
			sendResponseWithSignature(t, w, http.StatusOK, getDeviceServerResponse(t))
		case "session-server":
			sendResponseWithSignature(t, w, http.StatusOK, getSessionServerResponse(t))
		case "user-person/6084":
			switch r.Method {
			case http.MethodGet:
				sendResponseWithSignature(t, w, http.StatusOK, getUserPersonGetResponse(t))
			case http.MethodPut:
				sendResponseWithSignature(t, w, http.StatusOK, getGenericIDResponse(t))
			default:
				t.Errorf(errorRequestToUnMockedHTTPMethod, r.URL, r.Method)
			}
		case "user/6084/monetary-account/9512/draft-payment":
			switch r.Method {
			case http.MethodPost:
				sendResponseWithSignature(t, w, http.StatusOK, getGenericIDResponse(t))
			default:
				t.Errorf(errorRequestToUnMockedHTTPMethod, r.URL, r.Method)
			}
		case "user/6084/monetary-account-bank", "user/6084/monetary-account-bank/9601":
			sendResponseWithSignature(t, w, http.StatusOK, getMonetaryAccountBankGet(t))
		case "user/6084/monetary-account-savings", "user/6084/monetary-account-savings/9601":
			sendResponseWithSignature(t, w, http.StatusOK, getMonetaryAccountSavings(t))
		case "user/6084/monetary-account/9618/draft-payment", "user/6084/monetary-account/9618/draft-payment/6292":
			switch r.Method {
			case http.MethodPost, http.MethodPut:
				sendResponseWithSignature(t, w, http.StatusOK, getGenericIDResponse(t))
			case http.MethodGet:
				sendResponseWithSignature(t, w, http.StatusOK, getDraftPaymentGet(t))
			default:
				t.Errorf(errorRequestToUnMockedHTTPMethod, r.URL, r.Method)
			}

		case "user/6084/monetary-account/9520/mastercard-action/324":
			sendResponseWithSignature(t, w, http.StatusOK, getMasterCardActionGet(t))
		case "user/6084/monetary-account/10111/payment", "user/7082/monetary-account/10111/payment", "user/6084/monetary-account/10111/payment/1":
			sendResponseWithSignature(t, w, http.StatusOK, getPaymentGet(t))
		case "attachment-public/f9a1a89a-fdc1-4de5-89d5-e477cccd22c4/content":
			sendResponseWithSignature(t, w, http.StatusOK, getPaymentGet(t))
		case "/v1/session/133912", "v1/session/133912", "session/133912":
			sendResponseWithSignature(t, w, http.StatusOK, getSessionServerResponse(t))
		default:
			t.Errorf("requst made to an un mocked url: %v", r.URL)
		}
	})
}

func createBunqFakeHandlerWithError(t *testing.T, endpointToError string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == endpointToError {
			sendResponseWithSignature(t, w, http.StatusTeapot, getErrorResponse(t))
		} else {
			createBunqFakeHandler(t)(w, r)
		}
	})
}

func createClientWithFakeServer(t *testing.T) (*Client, *httptest.Server, context.CancelFunc) {
	fakeServer := httptest.NewServer(createBunqFakeHandler(t))

	key, err := CreateNewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := NewClient(ctx, fmt.Sprintf("%s/v1/", fakeServer.URL), key, "")

	return c, fakeServer, cancel
}

func getInstallationResponse(t *testing.T) *responseInstallation {
	var obj responseInstallation
	res := createResponseStruct(t, formatFilePathByName("installation_response"), &obj)

	return res.(*responseInstallation)
}

func getDeviceServerResponse(t *testing.T) *responseDeviceServer {
	var obj responseDeviceServer
	res := createResponseStruct(t, formatFilePathByName("device_server_response"), &obj)

	return res.(*responseDeviceServer)
}

func getSessionServerResponse(t *testing.T) *responseSessionServer {
	var obj responseSessionServer
	res := createResponseStruct(t, formatFilePathByName("session_server_response"), &obj)

	return res.(*responseSessionServer)
}

func getUserPersonGetResponse(t *testing.T) *responseUserPerson {
	var obj responseUserPerson
	res := createResponseStruct(t, formatFilePathByName("user_person_get_response"), &obj)

	return res.(*responseUserPerson)
}

func getMonetaryAccountBankGet(t *testing.T) *ResponseMonetaryAccountBankGet {
	var obj ResponseMonetaryAccountBankGet
	res := createResponseStruct(t, formatFilePathByName("monetary_account_bank_listing_response"), &obj)

	return res.(*ResponseMonetaryAccountBankGet)
}

func getMonetaryAccountSavings(t *testing.T) *ResponseMonetaryAccountSavingGet {
	var obj ResponseMonetaryAccountSavingGet
	res := createResponseStruct(t, formatFilePathByName("monetary_account_savings_response_get"), &obj)

	return res.(*ResponseMonetaryAccountSavingGet)
}

func getGenericIDResponse(t *testing.T) *responseBunqID {
	var obj responseBunqID
	res := createResponseStruct(t, formatFilePathByName("generic_id_response"), &obj)

	return res.(*responseBunqID)
}

func getDraftPaymentGet(t *testing.T) *responseDraftPaymentGet {
	var obj responseDraftPaymentGet
	res := createResponseStruct(t, formatFilePathByName("draft_payment_get_response"), &obj)

	return res.(*responseDraftPaymentGet)
}

func getMasterCardActionGet(t *testing.T) *responseMasterCardActionGet {
	var obj responseMasterCardActionGet
	res := createResponseStruct(t, formatFilePathByName("master_card_action_get_response"), &obj)

	return res.(*responseMasterCardActionGet)
}

func getPaymentGet(t *testing.T) *ResponsePaymentGet {
	var obj ResponsePaymentGet
	res := createResponseStruct(t, formatFilePathByName("payment_get_response"), &obj)

	return res.(*ResponsePaymentGet)
}

func getErrorResponse(t *testing.T) *responseError {
	var obj responseError
	res := createResponseStruct(t, formatFilePathByName("error_response"), &obj)

	return res.(*responseError)
}

func createResponseStruct(t *testing.T, path string, obj interface{}) interface{} {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(obj)
	if err != nil {
		t.Fatal(err)
	}

	return obj
}

func formatFilePathByName(fileName string) string {
	return fmt.Sprintf("../testdata/bunq/%s.json", fileName)
}

func sendResponseWithSignature(t *testing.T, w http.ResponseWriter, resCode int, body interface{}) {
	w.Header().Set("X-Bunq-Client-Response-Id", uuid.NewV4().String())
	w.Header().Set("Cache-controll", fmt.Sprintf("max-age=%s", uuid.NewV4().String()))
	stringToSign := fmt.Sprintf("%d\n", resCode)
	stringToSign += getAllHeaderToSignAsString(w.Header(), false)

	b, _ := json.Marshal(body)

	stringToSign += fmt.Sprintf("\n%s", b)

	h := sha256.New()
	_, _ = h.Write([]byte(stringToSign))

	signature, _ := rsa.SignPKCS1v15(rand.Reader, loadPrivateKey(), crypto.SHA256, h.Sum(nil))
	w.Header().Set("X-Bunq-Server-Signature", base64.StdEncoding.EncodeToString(signature))
	w.WriteHeader(resCode)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
}

var loadPrivateKeyOnce sync.Once
var privateKey *rsa.PrivateKey

func loadPrivateKey() *rsa.PrivateKey {
	loadPrivateKeyOnce.Do(func() {
		k, _ := ioutil.ReadFile("../testdata/bunq/private.key")

		block, _ := pem.Decode(k)
		privateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	})

	return privateKey
}

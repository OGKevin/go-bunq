package bunq

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	headerCacheControl     string = "Cache-Control"
	headerUserAgent        string = "User-Agent"
	headerXBunqLan         string = "X-Bunq-Language"
	headerXBunqRegion      string = "X-Bunq-Region"
	headerXBunqRequestID   string = "X-Bunq-Client-Request-Id"
	headerXBunqGeoLocation string = "X-Bunq-Geolocation"

	// BaseURLSandbox The base URL for the sanbox API.
	BaseURLSandbox string = "https://public-api.sandbox.bunq.com/v1/"
	// BaseURLProduction The base URL for the prod api
	BaseURLProduction string = "https://api.bunq.com/v1/"
)

// DetermineBaseURL returns the proper base url.
// If the env variable BUNQ_SANDBOX=true exists, sandbox is used, else production
func DetermineBaseURL() string {
	env := os.Getenv("BUNQ_SANDBOX")

	if env == "true" {
		return BaseURLSandbox
	}

	return BaseURLProduction
}

type queueEntry struct {
	req     *http.Request
	resChan chan *http.Response
	errChan chan error
}

type service struct {
	client *Client
}

// Client A client that can be used to communicate with the bunq api.
type Client struct {
	*http.Client
	ctx context.Context

	baseURL     string
	apiKey      string
	Debug       bool
	description string

	Err error

	requestQueue             chan queueEntry
	requestRateLimitMapMutex sync.RWMutex
	requestRateLimitMap      map[string]time.Time

	privateKey      *rsa.PrivateKey
	serverPublicKey *rsa.PublicKey

	userType

	// initOnce makes sure that init is called only once per instance. This will help that no new
	// new device keeps being registered etc.
	initOnce sync.Once

	tokenMutex sync.RWMutex
	// token is the token that needs to be in the auth header.
	token                *string
	installationContext  *installation
	sessionServerContext *sessionServer

	common                  service
	installation            *installationService
	deviceServer            *deviceServerService
	sessionServer           *sessionServerService
	UserService             *userService
	AccountService          *accountService
	PaymentService          *paymentService
	ScheduledPaymentService *scheduledPaymentService
	CardService             *cardService
	ContentService          *contentService
}

// NewClientFromContext create a new bunq client from a saved client context.
func NewClientFromContext(ctx context.Context, clientCtx *ClientContext) (*Client, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(clientCtx.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not parse private key")
	}

	block, _ := pem.Decode([]byte(clientCtx.InstallationContext.ServerPublicKey.ServerPublicKey))
	parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not parse server public key")
	}

	serverPubKey := parseResult.(*rsa.PublicKey)

	c := NewClient(ctx, clientCtx.BaseURL, privateKey, clientCtx.APIKey, "")
	c.apiKey = clientCtx.APIKey
	c.baseURL = clientCtx.BaseURL

	c.serverPublicKey = serverPubKey

	c.installationContext = clientCtx.InstallationContext
	c.sessionServerContext = clientCtx.SessionServerContext
	c.token = &c.sessionServerContext.Token.Token

	c.updateUserFlag()

	return c, nil
}

// NewClient create a new bunq client to use.
func NewClient(ctx context.Context, url string, key *rsa.PrivateKey, apikey, description string) *Client {
	c := Client{}
	c.ctx = ctx
	c.Client = http.DefaultClient
	c.baseURL = url
	c.description = description

	c.apiKey = apikey
	c.privateKey = key

	c.registerServices()

	return &c
}

// NewEmptyClient creates a new empty client.
func NewEmptyClient(ctx context.Context) *Client {
	c := Client{}
	c.ctx = ctx
	c.Client = http.DefaultClient
	c.baseURL = DetermineBaseURL()

	c.registerServices()

	return &c
}

func (c *Client) registerServices() {
	c.requestQueue = make(chan queueEntry, 9)
	c.requestRateLimitMap = make(map[string]time.Time)

	c.common.client = c

	c.installation = (*installationService)(&c.common)
	c.deviceServer = (*deviceServerService)(&c.common)
	c.sessionServer = (*sessionServerService)(&c.common)
	c.UserService = (*userService)(&c.common)
	c.PaymentService = (*paymentService)(&c.common)
	c.ScheduledPaymentService = (*scheduledPaymentService)(&c.common)
	c.AccountService = (*accountService)(&c.common)
	c.CardService = (*cardService)(&c.common)
	c.ContentService = (*contentService)(&c.common)

	c.spawnRequestHandlerWorker()
}

// SetAPIKey sets the api key
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// SetPrivateKey sets the api key
func (c *Client) SetPrivateKey(key *rsa.PrivateKey) {
	c.privateKey = key
}

// spawnRequestHandlerWorker will spawn a request queue worker that ensuers that all request
// that this client is making will be within the 1 request per second
// rate limit that bunq has.
//
// It starts when a new client has been created and dies when the client dies.
func (c *Client) spawnRequestHandlerWorker() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				close(c.requestQueue)
				return
			default:
				entry := <-c.requestQueue

				lastRequest := c.getLastExecutionTimeForRequest(entry.req)
				diff := time.Now().UTC().Sub(lastRequest)

				if diff.Seconds() < 1.0 {
					if c.Debug {
						log.Printf("bunq: waiting %f seconds before sending the http request.", 1.0-diff.Seconds())
					}

					time.Sleep(time.Duration((1.0 - diff.Seconds()) * float64(time.Second)))
				}

				go c.registerRequestInRateLimitMap(entry.req)

				if c.Debug {
					dump, _ := httputil.DumpRequest(entry.req, true)
					log.Printf("\n%s\n", dump)
				}

				res, err := c.Do(entry.req)

				if err != nil && c.Debug {
					log.Print(err)
				}

				if c.Debug && err == nil {
					dump, _ := httputil.DumpResponse(res, true)
					log.Printf("\n%s\n", dump)
				}

				entry.resChan <- res
				entry.errChan <- errors.Wrap(err, "bunq: http request failed.")
			}
		}
	}()
}

func (c *Client) getLastExecutionTimeForRequest(r *http.Request) time.Time {
	c.requestRateLimitMapMutex.RLock()
	defer c.requestRateLimitMapMutex.RUnlock()

	if t, ok := c.requestRateLimitMap[r.URL.Path]; ok {
		return t
	}

	return time.Now().UTC().Add(time.Duration(-2) * time.Second)
}

func (c *Client) registerRequestInRateLimitMap(r *http.Request) {
	c.requestRateLimitMapMutex.Lock()
	defer c.requestRateLimitMapMutex.Unlock()

	c.requestRateLimitMap[r.URL.Path] = time.Now().UTC()
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	err := c.setAllNeededHeader(r)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not set all required headers")
	}

	resChan := make(chan *http.Response, 1)
	errChan := make(chan error, 1)

	c.requestQueue <- queueEntry{
		req:     r,
		resChan: resChan,
		errChan: errChan,
	}

	res, err := <-resChan, <-errChan
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusInternalServerError {
			return nil, errors.New(fmt.Sprintf("bunq: http request failed with status %d", res.StatusCode))
		}

		errResponse := createErrorResponse(res)

		return nil, fmt.Errorf(
			"bunq: http request failed with status %d and description %q and response header: %q",
			res.StatusCode,
			errResponse.Error[0].ErrorDescription,
			res.Header.Get("X-Bunq-Client-Response-Id"),
		)
	}

	err = c.verifyResponse(r, res)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request was successful but repose verification failed")
	}

	return res, err
}

func (c *Client) setAllNeededHeader(r *http.Request) error {
	c.setAllDefaultHeader(r)
	var err error

	if shouldSignOrVerify(r.URL.Path) {
		c.tokenMutex.RLock()
		r.Header.Set("X-Bunq-Client-Authentication", *c.token)
		c.tokenMutex.RUnlock()

		err = c.addSignatureHeader(r)
	}

	return err
}

func shouldSignOrVerify(url string) bool {
	switch url[4:] {
	case "installation":
		return false
	default:
		return true
	}
}

func (*Client) setAllDefaultHeader(r *http.Request) {
	r.Header.Set(headerCacheControl, "no-cache")
	r.Header.Set(headerUserAgent, fmt.Sprintf("OGKevin-go-bunq-%s", os.Getenv("v0.1.2")))
	r.Header.Set(headerXBunqLan, "en_US")
	r.Header.Set(headerXBunqRegion, "nl_NL")
	r.Header.Set(headerXBunqGeoLocation, "0 0 0 0 NL")
	r.Header.Set(headerXBunqRequestID, generateRequestID())
}

func (c *Client) verifyResponse(r *http.Request, res *http.Response) error {
	if shouldSignOrVerify(r.URL.Path) {
		verified, err := c.verifySignature(res)

		if !verified {
			return errors.Wrap(err, "cloud not validate that request came from bunq")
		}
	}

	return nil
}

func (c *Client) formatRequestURL(path string) string {
	return c.baseURL + path
}

func createErrorResponse(r *http.Response) responseError {
	defer r.Body.Close()
	resBody, _ := ioutil.ReadAll(r.Body)

	var errorResponse responseError
	_ = json.Unmarshal(resBody, &errorResponse)

	return errorResponse
}

func generateRequestID() string {
	uid := uuid.NewV4()

	return uid.String()
}

// ExportClientContext exports the client context of the current client.
func (c *Client) ExportClientContext() (ClientContext, error) {
	p := x509.MarshalPKCS1PrivateKey(c.privateKey)
	userID, err := c.GetUserID()
	if err != nil {
		return ClientContext{}, err
	}

	ctx := ClientContext{
		PrivateKey:           p,
		InstallationContext:  c.installationContext,
		SessionServerContext: c.sessionServerContext,
		APIKey:               c.apiKey,
		BaseURL:              c.baseURL,
		UserID:               uint(userID),
	}

	return ctx, nil
}

// Init init's the client by preforming installation, device and session server where needed.
// this is a heavy task and should only be called once per context.
func (c *Client) Init() error {
	if c.Debug {
		log.Print("bunq: init client")
	}

	errChan := make(chan error, 1)

	c.initOnce.Do(func() {
		if c.installationContext == nil {
			c.preformNewInstallation(errChan)
		} else {
			if c.Debug {
				log.Print("bunq: installation context is not nil, only creating new session")
			}
			c.setInstallationToken()
			_, err := c.sessionServer.create()
			if err != nil {
				errChan <- errors.Wrap(err, "bunq: could not create new session")
				return
			}
		}

		c.spawnSessionHandlingWorker()
	})

	if len(errChan) != 0 {
		return <-errChan
	}

	return nil
}

func (c *Client) preformNewInstallation(errChan chan error) {
	if c.Debug {
		log.Print("bunq: installation context is nil, doing installation, device-server and session-server calls")
	}

	_, err := c.installation.create()
	if err != nil {
		errChan <- errors.Wrap(err, "bunq: could not init installation")
		return
	}

	_, err = c.deviceServer.create()
	if err != nil {
		errChan <- errors.Wrap(err, "bunq: could not init device server")
		return
	}

	_, err = c.sessionServer.create()
	if err != nil {
		errChan <- errors.Wrap(err, "bunq: could not init session server")
		return
	}
}

// IsUserPerson returns true if the current auth user is of type UserPerson
func (c *Client) IsUserPerson() bool {
	return c.isUserPerson
}

// IsUserCompany returns true if the current auth user is of type UserCompany
func (c *Client) IsUserCompany() bool {
	return c.isUserCompany
}

// IsUserAPIKey returns true if the current auth user is of type UserApiKey
func (c *Client) IsUserAPIKey() bool {
	return c.isUserAPIkey
}

func (c *Client) updateUserFlag() {
	if c.sessionServerContext.UserPerson.ID != 0 {
		c.isUserPerson = true
	} else if c.sessionServerContext.UserCompany.ID != 0 {
		c.isUserCompany = true
	} else if c.sessionServerContext.UserAPIKey.ID != 0 {
		c.isUserAPIkey = true
	}
}

// spawnSessionHandlingWorker makes sure that the user session is always valid. This is to ensure that no 403
// errors happen. The session is valid based on the user's auto logout time in the bunq app.
func (c *Client) spawnSessionHandlingWorker() {
	go func() {
		if c.Debug {
			log.Print("bunq: spawned session handling worker")
		}

		for {
			select {
			case <-c.ctx.Done():
				err := c.sessionServer.delete()
				if err != nil {
					c.Err = errors.Wrap(err, "bunq: session handler: could not delete session")
				}
				return
			default:
				expSec, err := c.getSessionExpInSec()
				if err != nil {
					c.Err = errors.Wrap(err, "bunq: could not get exp time")
					continue
				}

				expTime := time.Now().UTC().Add(time.Second * time.Duration(expSec-5))

				if c.Debug {
					log.Printf("bunq: session will expirte at %q", expTime)
				}

				if expTime.After(time.Now().UTC()) {
					timeToSleep := expTime.Sub(time.Now().UTC())

					if c.Debug {
						log.Printf("bunq: session worker will sleep for %f seconds until it renews the session.", timeToSleep.Seconds())
					}

					time.Sleep(timeToSleep)
				}

				c.setInstallationToken()
				_, err = c.sessionServer.create()
				if err != nil {
					c.Err = errors.Wrap(err, "bunq: session handler: could not create session")
				}
			}
		}
	}()
}

func (c *Client) getSessionExpInSec() (int64, error) {
	if c.IsUserPerson() {
		return c.sessionServerContext.UserPerson.SessionTimeout, nil
	} else if c.IsUserCompany() {
		return c.sessionServerContext.UserCompany.SessionTimeout, nil
	}

	return 0, fmt.Errorf("bunq: could not get user expirty time")
}

func (c *Client) setInstallationToken() {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	c.token = &c.installationContext.Token.Token
}

// GetUserID returns the user id of the current auth user.
func (c *Client) GetUserID() (int, error) {
	if c.isUserPerson {
		return c.sessionServerContext.UserPerson.ID, nil
	} else if c.isUserCompany {
		return c.sessionServerContext.UserCompany.ID, nil
	} else if c.isUserAPIkey {
		return c.sessionServerContext.UserAPIKey.ID, nil
	}

	return 0, fmt.Errorf("bunq: could not determine user id")
}

func (c *Client) preformRequest(method, url string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("bunq: could not create request for  %s", url))
	}

	res, err := c.do(r)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("bunq: request to %s failed", url))
	}

	return res, err
}

func (c *Client) parseResponse(res *http.Response, obj interface{}) error {
	defer res.Body.Close()

	err := json.NewDecoder(res.Body).Decode(obj)
	if err != nil {
		return errors.Wrap(err, "bunq: could not parse response")
	}

	return nil
}

func (c *Client) doCURequest(url string, bodyRaw []byte, httpMethod string) (*responseBunqID, error) {
	res, err := c.preformRequest(httpMethod, url, bytes.NewBuffer(bodyRaw))
	if err != nil {
		return nil, err
	}

	var resBunqID responseBunqID

	return &resBunqID, c.parseResponse(res, &resBunqID)
}

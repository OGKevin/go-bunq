package bunq

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePaymentService_CreateBatchPayment() {
	key, err := CreateNewKeyPair()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := NewClient(ctx, BaseURLSandbox, key, "sandbox_ab7df7985a66133b1abecf42871801edaafe5bc51ef9769f5a032876", "My awesome app")
	err = c.Init()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		log.Print(i)

		_, err = c.PaymentService.CreatePaymentBatch(
			10111,
			PaymentBatchCreate{
				Payments: generateBatchEntries(100),
			},
		)
		if err != nil {
			panic(err)
		}
	}
}

func TestDraftPaymentCreate(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := createNewDraftPayment(c)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].ID.ID)
}

func TestGetDraftPayment(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := createNewDraftPayment(c)

	assert.NoError(t, err)

	resGet, err := c.PaymentService.GetDraftPayment(res.Response[0].ID.ID, 9618)

	assert.NoError(t, err)
	assert.NotZero(t, resGet.Response[0].DraftPayment.ID)
	assert.NotEmpty(t, resGet.Response[0].DraftPayment.Updated)
}

func TestDraftPaymentUpdate(t *testing.T) {
	t.Parallel()

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	assert.NoError(t, c.Init())

	res, err := createNewDraftPayment(c)

	assert.NoError(t, err)
	assert.NotZero(t, res.Response[0].ID.ID)

	resGet, err := c.PaymentService.GetDraftPayment(res.Response[0].ID.ID, 9618)

	assert.NoError(t, err)
	assert.NotZero(t, resGet.Response[0].DraftPayment.Updated)

	allDraftPaymentEntry := convertDraftPaymentEntryToCreateEntry(resGet.Response[0].DraftPayment.Entries...)

	_, err = c.PaymentService.UpdateDraftPayment(
		res.Response[0].ID.ID,
		9618,
		requestUpdateDraftPayment{
			requestCreateDraftPayment: requestCreateDraftPayment{
				Entries: append(allDraftPaymentEntry, draftPaymentEntryCreate{
					Amount: Amount{
						Currency: "EUR",
						Value:    "1",
					},
					CounterpartyAlias: Pointer{
						PType: "EMAIL",
						Value: "bravo@bunq.com",
					},
					Description: "test",
				}),
			},
			UpdatedTimestamp: resGet.Response[0].DraftPayment.Updated,
		},
	)

	assert.NoError(t, err)

	_, err = c.PaymentService.GetDraftPayment(res.Response[0].ID.ID, 9618)
	assert.NoError(t, err)
}

func createNewDraftPayment(c *Client) (*responseBunqID, error) {
	i := 1
	return c.PaymentService.CreateDraftPayment(
		9618,
		requestCreateDraftPayment{
			Entries: []draftPaymentEntryCreate{
				{
					Amount: Amount{
						Currency: "EUR",
						Value:    "1",
					},
					CounterpartyAlias: Pointer{
						PType: "EMAIL",
						Value: "bravo@bunq.com",
					},
					Description: "test",
				},
			},
			NumberOfRequiredAccepts: &i,
		},
	)
}

func generateBatchEntries(nr int) []PaymentCreate {
	var entries []PaymentCreate

	for i := 0; i < nr; i++ {
		entries = append(
			entries,
			PaymentCreate{
				Amount: Amount{
					Currency: "EUR",
					Value:    "0.01",
				},
				CounterpartyAlias: Pointer{
					PType: "EMAIL",
					Value: "bravo@bunq.com",
				},
				Description: "test",
			},
		)
	}

	return entries
}

func convertDraftPaymentEntryToCreateEntry(allEntry ...draftPaymentEntry) []draftPaymentEntryCreate {
	var allCreateEntry []draftPaymentEntryCreate

	for _, entry := range allEntry {
		allCreateEntry = append(allCreateEntry, draftPaymentEntryCreate{
			Amount: entry.Amount,
			CounterpartyAlias: Pointer{
				PType: "IBAN",
				Value: entry.CounterpartyAlias.IBAN,
				Name:  &entry.CounterpartyAlias.DisplayName,
			},
			Description: entry.Description,
		})
	}

	return allCreateEntry
}

func Test_paymentService_GetAllPayment(t *testing.T) {
	t.Parallel()

	type fields struct {
		client *Client
	}
	type args struct {
		monetaryAccountID uint
	}

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "main",
			fields: fields{client: c},
			args:   args{monetaryAccountID: 10111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.fields.client.Init())

			p := &paymentService{
				client: tt.fields.client,
			}
			got, err := p.GetAllPayment(tt.args.monetaryAccountID)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}

			assert.NotZero(t, got.Response[0].Payment.ID)
		})
	}
}

func Test_paymentService_GetAllOlderPayment(t *testing.T) {
	t.Parallel()

	type fields struct {
		client *Client
	}

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()
	if !assert.NoError(t, c.Init()) {
		return
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "main",
			fields: fields{
				client: c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := c.PaymentService.GetAllPayment(10111)
			if !assert.NoError(t, err) {
				return
			}

			p := &paymentService{
				client: tt.fields.client,
			}
			got, err := p.GetAllOlderPayment(res.Pagination)

			if assert.NoError(t, err) {
				assert.NotZero(t, got.Response[0].Payment.ID)
			}
		})
	}
}

func Test_paymentService_GetPayment(t *testing.T) {
	t.Parallel()

	type fields struct {
		client *Client
	}

	c, fakeServer, cancel := createClientWithFakeServer(t)
	defer cancel()
	defer fakeServer.Close()
	if !assert.NoError(t, c.Init()) {
		return
	}

	type args struct {
		monetaryAccountID uint
		paymentID         uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "main",
			fields: fields{
				client: c,
			},
			args: args{
				monetaryAccountID: 10111,
				paymentID:         1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &paymentService{
				client: tt.fields.client,
			}
			got, err := p.GetPayment(tt.args.monetaryAccountID, tt.args.paymentID)

			if assert.NoError(t, err) {
				assert.NotZero(t, got.Response[0].Payment.ID)
			}
		})
	}
}

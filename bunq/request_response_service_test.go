package bunq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_requestResponseService_GetAllRequestResponses(t *testing.T) {
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
			args:   args{monetaryAccountID: 9999},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.fields.client.Init())

			rs := &requestResponseService{
				client: tt.fields.client,
			}
			got, err := rs.GetAllRequestResponses(tt.args.monetaryAccountID)
			if (err != nil) != tt.wantErr {
				assert.NoError(t, err)
				return
			}

			assert.NotZero(t, got.Response[0].RequestResponse.ID)
		})
	}
}

func Test_requestResponseService_GetAllOlderRequestResponses(t *testing.T) {
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
			res, err := c.RequestResponseService.GetAllRequestResponses(9999)
			if !assert.NoError(t, err) {
				return
			}

			rs := &requestResponseService{
				client: tt.fields.client,
			}
			got, err := rs.GetAllOlderRequestResponses(res.Pagination)

			if assert.NoError(t, err) {
				assert.NotZero(t, got.Response[0].RequestResponse.ID)
			}
		})
	}
}

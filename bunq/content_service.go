package bunq

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

type contentService service

func (c *contentService) GetAttachmentPublic(id string) (string, error) {
	res, err := c.client.preformRequest(http.MethodGet, c.client.formatRequestURL(fmt.Sprintf("attachment-public/%s/content", id)), nil)
	if err != nil {
		return "", errors.Wrap(err, "bunq: request to get attachment content failed")
	}

	pr, pw := io.Pipe()
	encoder := base64.NewEncoder(base64.StdEncoding, pw)

	go func() {
		_, err := io.Copy(encoder, res.Body)
		defer encoder.Close()
		defer pw.CloseWithError(err)
	}()

	out, err := ioutil.ReadAll(pr)
	if err != nil {
		return "", errors.Wrap(err, "bunq: could not read from reader")
	}

	return string(out), nil
}

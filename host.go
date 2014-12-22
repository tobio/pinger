package pinger

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Host is a representation of a host to ping,
// it contains optional ExpectedBody and ExpectedStatusCode fields
type Host struct {
	Name               string `json:"name"`
	Url                string `json:"url"`
	ExpectedBody       string `json:"expected_body"`
	ExpectedStatusCode int    `json:"expected_statusCode"`
}

// Error variables
var (
	ErrBodyMismatch       = errors.New("body mismatch")
	ErrStatusCodeMismatch = errors.New("status code mismatch")
	ErrBadStatusCode      = errors.New("status code is not 200")
)

// Ping pings the host,
// it returns the response HTTP status code, body and error
func (h Host) Ping(c *http.Client) (int, []byte, error) {
	res, err := c.Get(h.Url)
	if err != nil {
		return 0, nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	// compare the body
	if h.ExpectedBody != "" && string(b) != h.ExpectedBody {
		err = ErrBodyMismatch
	}

	// compare the status code
	if h.ExpectedStatusCode != 0 && h.ExpectedStatusCode != res.StatusCode {
		err = ErrStatusCodeMismatch
	} else if h.ExpectedStatusCode == 0 && res.StatusCode != 200 {
		err = ErrBadStatusCode
	}

	return res.StatusCode, b, err
}

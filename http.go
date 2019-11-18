package tado

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const defaultBaseURL = "https://my.tado.com/api"

func (c *Client) do(in input, out interface{}) error {
	// ensure accesstoken is still valid
	err := c.validateAccessToken()
	if err != nil {
		return err
	}

	// encode input as JSON if needed
	var body io.Reader
	switch in.method() {
	case http.MethodPost, http.MethodPut:
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return fmt.Errorf("error encoding input: %s", err)
		}
		body = buf
	}

	// create HTTP request
	req, err := http.NewRequest(in.method(), c.baseURL+in.path(), body)
	if err != nil {
		return err
	}

	// set authentication header
	req.Header.Set("Authorization", "Bearer "+c.tr.AccessToken)

	// execute HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP error: %s", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// check HTTP status
	if resp.StatusCode >= http.StatusBadRequest {
		// not OK, read body
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<14))
		if err != nil {
			return fmt.Errorf("HTTP error: %s", err)
		}
		// return the body as error
		return fmt.Errorf("error: HTTP status %d: %s", resp.StatusCode, string(body))
	}

	// OK response, decode into output
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return fmt.Errorf("error decoding output: %s", err)
	}

	return nil
}

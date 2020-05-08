package corezoid

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (c *Client) Task(task Task) (*Op, error) {
	converted, err := c.convert(task)
	if err != nil {
		return nil, err
	}

	return c.CallOne(converted, true)
}

func (c *Client) AsyncTask(task Task) (*Op, error) {
	converted, err := c.convert(task)
	if err != nil {
		return nil, err
	}

	return c.CallOne(converted, false)
}

func (c *Client) CallOne(op Op, synchronous bool) (*Op, error) {
	if res, err := c.Call(Ops{List: []Op{op}}, synchronous); err != nil {
		return nil, err
	} else {
		if !res.IsRequestOK() {
			return nil, fmt.Errorf("got non-ok status: %s and http status code %d", res.RequestProc, res.StatusCode)
		}
		if len(res.List) != 1 {
			return nil, fmt.Errorf("got %d ops but 1 was expected", len(res.List))
		}

		x := res.List[0]

		return &x, nil
	}

}

func (c *Client) Call(ops Ops, synchronous bool) (*OpsResult, error) {
	// Target host
	var endpoint string
	if synchronous {
		endpoint = syncApi
	} else {
		endpoint = asyncApi
	}

	// Content encoding
	content, contentType, err := c.encode(ops)
	if err != nil {
		return nil, err
	}

	// Content path type (exp "json")
	var pathType string
	switch contentType {
	case "application/json":
		pathType = "json"
		break
	}

	now := time.Now().Second()

	// Sign
	signedPath, err := c.encrypt([]byte(fmt.Sprintf("%d%s%s%s", now, c.apiSecret, content, c.apiSecret)))
	if err != nil {
		return nil, err
	}

	// Building request
	uri := fmt.Sprintf("%s/%s/%s/%d/%s", endpoint, pathType, c.apiKey, now, signedPath)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	// Sending request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result *OpsResult

	if resp.ContentLength > 0 {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read response body - %s", err)
		}

		result, err = c.decode(bts)
		if err != nil {
			return nil, fmt.Errorf("unable to decode response body - %s", err)
		}
	}

	if result != nil {
		result.StatusCode = resp.StatusCode
	} else {
		return nil, fmt.Errorf("got %d HTTP status code from Corezoid", resp.StatusCode)
	}

	return result, nil
}

func (c *Client) encrypt(source []byte) ([]byte, error) {
	h := sha1.New()
	if _, err := h.Write(source); err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(h.Sum(nil))), nil
}

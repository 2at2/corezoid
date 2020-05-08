package corezoid

import (
	"encoding/json"
	"errors"
)

func (c *Client) convert(source interface{}) (Op, error) {
	var op Op

	bts, err := json.Marshal(source)
	if err != nil {
		return op, err
	}
	if err := json.Unmarshal(bts, &op); err != nil {
		return op, err
	}

	return op, nil
}

func (c *Client) encode(ops Ops) ([]byte, string, error) {
	switch c.contentType {
	case Json:
		t := "application/json"
		payload, err := json.Marshal(ops)
		if err != nil {
			return nil, t, err
		}

		return payload, t, nil

	default:
		return nil, "", errors.New("unknown content type")
	}
}

func (c *Client) decode(data []byte) (*OpsResult, error) {
	var res OpsResult

	switch c.contentType {
	case Json:
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, err
		}
		break

	default:
		return nil, errors.New("unknown content type")
	}

	return &res, nil
}

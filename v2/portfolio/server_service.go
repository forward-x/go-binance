package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/ping",
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	d := map[string]string{}
	err = json.Unmarshal(data, &d)
	return err
}

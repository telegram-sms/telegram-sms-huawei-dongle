package client

import "net/http"

type RequestOptions interface {
	BeforeRequest(req *http.Request)
	TransformBody(c *Client, body []byte) []byte
}

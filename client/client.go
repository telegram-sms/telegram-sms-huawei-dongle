package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client/cookiejar"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client/fifo"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL string
	Host    string
	Tokens  *fifo.TokenQueue

	// The client with cookie jar
	client http.Client

	// Specify a path to update session cookie.
	// default to "/html/footer.html" (a very small file that assigns session cookie)
	SessionPath string
}

func (c *Client) url(path string) string {
	return c.BaseURL + path
}

var xmlHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>`)

func (c *Client) Request(path string, body []byte, opt RequestOptions) ([]byte, error) {
	token := ""
	action := "GET"
	if body != nil {
		action = "POST"

		body = append(xmlHeader, body...)
		token = c.Tokens.Consume()

		if opt != nil {
			body = opt.TransformBody(c, body)
		}
	}

	req, err := http.NewRequest(action, c.url(path), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Host", c.Host)
	req.Header.Set("Origin", c.BaseURL)
	if token != "" {
		req.Header["__RequestVerificationToken"] = []string{token}
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", c.url("/html/home.html"))
	if opt != nil {
		opt.BeforeRequest(req)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %w", err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = resp.Body.Close()
	}

	c.Tokens.Add(resp.Header.Get("__RequestVerificationTokenOne"))
	c.Tokens.Add(resp.Header.Get("__RequestVerificationTokenTwo"))
	c.Tokens.Add(resp.Header.Get("__RequestVerificationToken"))

	return result, err
}

func (c *Client) API(path string, body interface{}, resp interface{}, opt RequestOptions) error {
	//log.Printf("xml.req: %s: %s\n", path, body)
	bodyBytes, err := xml.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not convert request body to xml: %w", err)
	}

	respBytes, err := c.Request("/api"+path, bodyBytes, opt)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	//log.Printf("xml.resp: %s\n", respBytes)

	return parseResp(respBytes, resp)
}

func (c *Client) Init(baseURL string) error {
	uri, _ := url.Parse(baseURL)
	return c.InitWithHost(baseURL, uri.Host)
}

func (c *Client) InitWithHost(baseURL, host string) error {
	if baseURL[len(baseURL)-1] == '/' {
		c.BaseURL = baseURL[:len(baseURL)-1]
	} else {
		c.BaseURL = baseURL
	}

	if c.SessionPath == "" {
		c.SessionPath = "/html/footer.html"
	}

	c.Host = host

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	c.client = http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	c.Tokens = &fifo.TokenQueue{}
	c.Tokens.Init()

	return nil
}

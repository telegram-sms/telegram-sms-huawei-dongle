package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	BaseURL string
	client  http.Client
	Host    string
	Token   string
}

func (c *Client) url(path string) string {
	return c.BaseURL + path
}

var xmlHeader = []byte(`<?xml version="1.0" encoding="UTF-8"?>`)

func (c *Client) Request(path string, body []byte, opt RequestOptions) ([]byte, error) {
	action := "GET"
	if body != nil {
		action = "POST"

		body = append(xmlHeader, body...)

		if opt != nil {
			body = opt.TransformBody(c, body)
		}
	}

	req, err := http.NewRequest(action, c.url(path), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Host", c.Host)
	req.Header.Set("__RequestVerificationToken", c.Token)
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

	u, _ := url.Parse(c.BaseURL)
	for _, cookie := range c.client.Jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	return result, err
}

func (c *Client) API(path string, body interface{}, resp interface{}, opt RequestOptions) error {
	log.Printf("xml.req: %s: %s\n", path, body)
	bodyBytes, err := xml.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not convert request body to xml: %w", err)
	}

	respBytes, err := c.Request("/api"+path, bodyBytes, opt)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	log.Printf("xml.resp: %s\n", respBytes)

	err = xml.Unmarshal(respBytes, resp)
	return err
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

	return nil
}

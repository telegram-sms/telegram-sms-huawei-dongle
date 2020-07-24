package client

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) UpdateSession() error {
	_, err := c.Request("/", nil, nil)
	if err == nil {
		return nil
	}

	// Use API to fetch session id
	session, err := c.GetSessionTokenInfo()
	if err == nil {
		sessionID := strings.Replace(session.Session, "SessionID=", "", 1)
		path, err := url.Parse(c.BaseURL)
		if err != nil {
			return err
		}
		cookie := &http.Cookie{
			Name:    "SessionID",
			Value:   sessionID,
			Expires: time.Now().AddDate(99, 0, 0),
		}
		c.client.Jar.SetCookies(path, []*http.Cookie{cookie})
	}

	return nil
}

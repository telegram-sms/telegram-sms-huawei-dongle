package client

import (
	"net/http"
	"net/url"
	"strings"
)

// UpdateSession: Call this method periodically to keep session alive (I think)
func (c *Client) UpdateSession() error {
	_, err := c.Request(c.SessionPath, nil, nil)
	if err == nil {
		return nil
	}

	// Use API to fetch session id (fallback)
	session, err := c.GetSessionTokenInfo()
	if err == nil {
		sessionID := strings.Replace(session.Session, "SessionID=", "", 1)
		path, err := url.Parse(c.BaseURL)
		if err != nil {
			return err
		}
		cookie := &http.Cookie{
			Name:  "SessionID",
			Value: sessionID,
		}
		c.client.Jar.SetCookies(path, []*http.Cookie{cookie})
	}

	return nil
}

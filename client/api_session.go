package client

import (
	"log"
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
	return c.UpdateSessionUsingAPI()
}

func (c *Client) UpdateSessionUsingAPI() error {
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
		log.Printf("sessionID: %s\n", sessionID)
	}

	return nil
}

// +build integration

package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

var dongleURL = "http://192.168.4.1/"
var dongleURLObj, _ = url.Parse(dongleURL)

func assertSessionCookie(t *testing.T, client Client) {
	found := false
	for _, cookie := range client.client.Jar.Cookies(dongleURLObj) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
		found = found || (cookie.Name == "SessionID" && cookie.Value != "")
	}
	assert.Truef(t, found, "should find cookie value for SessionID")
}

func TestClient_UpdateSession(t *testing.T) {
	client := Client{}
	err := client.Init(dongleURL)
	assert.Nilf(t, err, "should be able to access home page")
	_, err = client.GetSessionTokenInfo()
	assert.Nilf(t, err, "should be able to get session id")
	assertSessionCookie(t, client)
}

func TestClient_UpdateSessionUsingAPI(t *testing.T) {
	client := Client{}
	err := client.Init(dongleURL)
	assert.Nilf(t, err, "should be able to access home page")
	_, err = client.GetSessionTokenInfo()
	assert.Nilf(t, err, "should be able to get session id")
	assertSessionCookie(t, client)
}

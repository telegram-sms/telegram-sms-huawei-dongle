// +build integration

package client

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func getClient(t *testing.T) *Client {
	c := &Client{}
	err := c.Init(dongleURL)
	assert.Nilf(t, err, "could not init")

	return c
}

func getAdminClient(t *testing.T) *Client {
	c := getClient(t)
	success, err := c.Login("admin", os.Getenv("SECRET_PASSWORD"))
	assert.Nilf(t, err, "could not login")
	assert.Truef(t, success.IsLoginSuccess(), "login should success")

	login, err := c.GetLoginState()
	assert.Nilf(t, err, "could not get login state")
	assert.Truef(t, login.IsLoggedIn(), "should be logged in state")

	return c
}

func TestClient_Login(t *testing.T) {
	_ = getAdminClient(t)
}

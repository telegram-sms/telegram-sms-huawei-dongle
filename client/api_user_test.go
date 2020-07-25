package client

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_Login(t *testing.T) {
	c := &Client{}
	err := c.Init(dongleURL)
	assert.Nilf(t, err, "could not init")
	success, err := c.Login("admin", os.Getenv("SECRET_PASSWORD"))
	assert.Nil(t, err, "could not login")
	assert.Truef(t, success, "login should success")

	login, err := c.GetLoginState()
	assert.Nilf(t, err, "could not get login state")
	assert.Truef(t, login.IsLoggedIn(), "should be logged in state")
}

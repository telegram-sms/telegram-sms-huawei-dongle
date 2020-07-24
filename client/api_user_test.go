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
	resp, err := c.Login("admin", os.Getenv("SECRET_PASSWORD"))
	_ = resp
	assert.Nilf(t, err, "could not login")
}

package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_Login(t *testing.T) {
	fmt.Printf("SECRET_PASSWORD: %s\n", os.Getenv("SECRET_PASSWORD"))

	c := &Client{}
	err := c.Init(dongleURL)
	assert.Nilf(t, err, "could not init")
	resp, err := c.Login("admin", os.Getenv("SECRET_PASSWORD"), false)
	_ = resp
	assert.Nilf(t, err, "could not login")
}

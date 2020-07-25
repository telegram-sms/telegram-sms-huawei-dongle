// +build integration

package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_SMSList(t *testing.T) {
	c := getAdminClient(t)
	resp, err := c.SMSList(1, 20)
	assert.Nilf(t, err, fmt.Sprintf("could not get messages: %s", err))
	msg := resp.Messages[0]
	fmt.Printf("[%s, %s]: %s", msg.Date, msg.Phone, msg.Content)
}

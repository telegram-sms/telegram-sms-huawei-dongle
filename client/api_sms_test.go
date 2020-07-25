// +build integration

package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_SMSList(t *testing.T) {
	c := getAdminClient(t)
	resp, err := c.SMSList(1, 20)
	assert.Nilf(t, err, fmt.Sprintf("could not get messages: %s", err))
	msg := resp.Messages[0]
	fmt.Printf("[%s, %s]: %s", msg.Date, msg.Phone, msg.Content)
}

func TestClient_SendSMS(t *testing.T) {
	c := getAdminClient(t)

	// Check if SMS received.
	_, _ = c.SendSMS(os.Getenv("SECRET_PHONE"), "This is a test :D")
}

func TestClient_GetSendStatus(t *testing.T) {
	c := getAdminClient(t)
	resp, err := c.GetSendStatus()
	assert.Nilf(t, err, "api call should success")
	fmt.Printf("Send status: %d/%d\n", resp.CurIndex, resp.TotalCount)
}

// +build integration

package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestClient_GetDeviceStatus(t *testing.T) {
	c := getAdminClient(t)

	status, err := c.GetDeviceStatus()
	assert.Nilf(t, err, "error should be nil")
	assert.Zerof(t, status.ErrorCode, "should success")
	fmt.Printf("has battery: %s\n", strconv.FormatBool(status.HasBattery()))
}

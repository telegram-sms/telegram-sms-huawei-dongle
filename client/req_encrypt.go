package client

import (
	"crypto/rsa"
	"fmt"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client/crypto"
	"net/http"
)

type EncryptedRequest struct {
	pubKey *rsa.PublicKey
}

func (e *EncryptedRequest) BeforeRequest(req *http.Request) {
	req.Header["encrypt_transmit"] = []string{"encrypt_transmit"}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
}

func (e *EncryptedRequest) TransformBody(_ *Client, body []byte) []byte {
	fmt.Printf("before: %s\n", string(body))
	encrypted := crypto.EncryptHuaweiRSA(body, e.pubKey)
	fmt.Printf("after:  %s\n", encrypted)
	return []byte(encrypted)
}

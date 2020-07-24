package client

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client/crypto"
	"net/http"
)

type EncryptedRequest struct {
	pubKey *rsa.PublicKey
}

func (e *EncryptedRequest) BeforeRequest(req *http.Request) {
	req.Header.Add("encrypt_transmit", "encrypt_transmit")
}

func (e *EncryptedRequest) TransformBody(_ *Client, body []byte) []byte {
	fmt.Printf("before: %s\n", string(body))
	encrypted := hex.EncodeToString(crypto.EncryptHuaweiRSA(body, e.pubKey))
	fmt.Printf("after:  %s\n", encrypted)
	return []byte(encrypted)
}

package client

import (
	"encoding/xml"
	"fmt"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client/crypto"
	"log"
)

type LoginStateResp struct {
	XMLName            xml.Name `xml:"response"`
	State              string   `xml:"State"`
	Username           string   `xml:"Username"`
	PasswordType       string   `xml:"password_type"`
	ExternPasswordType string   `xml:"extern_password_type"`
	FirstLogin         string   `xml:"firstlogin"`
}

func (r *LoginStateResp) UseScarmLogin() bool {
	return r.ExternPasswordType == "1"
}

func (r *LoginStateResp) IsPasswordSalted() bool {
	return !r.UseScarmLogin() && r.PasswordType == "4"
}

func (c *Client) GetLoginState() (*LoginStateResp, error) {
	resp := LoginStateResp{}
	err := c.API("/user/state-login", nil, &resp, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get login state: %w", err)
	}

	return &resp, err
}

type loginPayload struct {
	XMLName      xml.Name `xml:"request"`
	Username     string   `xml:"Username"`
	Password     string   `xml:"Password"`
	PasswordType string   `xml:"password_type"`
}

type LoginResp struct {
}

func (c *Client) Login(username, password string) (*LoginResp, error) {
	if err := c.UpdateSession(); err != nil {
		return nil, fmt.Errorf("could not renew session: %w", err)
	}

	login, err := c.GetLoginState()
	if err != nil {
		return nil, err
	}

	if login.UseScarmLogin() {
		return nil, fmt.Errorf("unsupported login type: SCARM")
	}

	if login.IsPasswordSalted() {
		sess, err := c.GetSessionTokenInfo()
		if err != nil {
			return nil, fmt.Errorf("could not fetch session token: %w", err)
		}
		log.Printf("using salted password (token: %s)\n", sess.Token)
		password = crypto.EncodeSaltedPassword(username, password, sess.Token)
	} else {
		log.Println("using base64 for password")
		password = crypto.B64(password)
	}

	switches, err := c.GetModuleSwitches()
	if err != nil {
		return nil, fmt.Errorf("could not get global module switches: %w", err)
	}
	var opts RequestOptions = nil
	if switches.IsEncryptionEnabled() {
		pubKey, err := c.GetPublicKey()
		if err != nil {
			return nil, fmt.Errorf("could not fetch rsa key: %w", err)
		}
		opts = &EncryptedRequest{pubKey: pubKey}
	}

	payload := loginPayload{
		Username:     username,
		Password:     password,
		PasswordType: login.PasswordType,
	}
	resp := LoginResp{}
	err = c.API("/user/login", payload, &resp, opts)
	return &resp, err
}

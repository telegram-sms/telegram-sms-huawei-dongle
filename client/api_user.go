package client

import (
	"encoding/xml"
	"fmt"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client/crypto"
	"log"
)

type LoginStateResp struct {
	BaseResp

	State              int    `xml:"State"`
	Username           string `xml:"Username"`
	PasswordType       int    `xml:"password_type"`
	ExternPasswordType int    `xml:"extern_password_type"`
	FirstLogin         int    `xml:"firstlogin"`
}

func (r *LoginStateResp) IsLoggedIn() bool {
	return r.State == 0
}

func (r *LoginStateResp) UseScarmLogin() bool {
	return r.ExternPasswordType == 1
}

func (r *LoginStateResp) IsPasswordSalted() bool {
	return !r.UseScarmLogin() && r.PasswordType == 4
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
	PasswordType int      `xml:"password_type"`
}

type LoginResp struct {
	BaseResp

	Value string `xml:",chardata"`
}

func (r *LoginResp) IsLoginSuccess() bool {
	return r.Value == "OK"
}

// Login performs the login routine.
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

	var opts RequestOptions = nil
	switches, err := c.GetModuleSwitches()
	if err != nil {
		return nil, fmt.Errorf("could not get global module switches: %w", err)
	}

	if switches.IsEncryptionEnabled() {
		pubKey, err := c.GetPublicKey()
		if err != nil {
			return nil, fmt.Errorf("could not fetch rsa key: %w", err)
		}
		opts = &EncryptedRequest{pubKey: pubKey}
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

	payload := loginPayload{
		Username:     username,
		Password:     password,
		PasswordType: login.PasswordType,
	}
	resp := &LoginResp{}
	err = c.API("/user/login", payload, resp, opts)

	if resp.Value == "OK" {
		resp.ErrorCode = 0
	}
	return resp, err
}

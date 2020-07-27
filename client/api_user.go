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

func (c *Client) encodePassword(username, password string, salted bool) (string, error) {
	encoded := ""

	if salted {
		sess, err := c.GetSessionTokenInfo()
		if err != nil {
			return "", fmt.Errorf("could not fetch session token: %w", err)
		}
		log.Printf("using salted password (token: %s)\n", sess.Token)
		encoded = crypto.EncodeSaltedPassword(username, password, sess.Token)
	} else {
		log.Println("using base64 for password")
		encoded = crypto.B64(password)
	}

	return encoded, nil
}

func (c *Client) doLogin(username, encodedPassword string, pwType int, opts RequestOptions) (*LoginResp, error) {
	payload := loginPayload{
		Username:     username,
		Password:     encodedPassword,
		PasswordType: pwType,
	}
	resp := &LoginResp{}
	err := c.API("/user/login", payload, resp, opts)

	if resp.Value == "OK" {
		resp.ErrorCode = 0
	}
	return resp, err
}

// Login performs the login routine.
// It ignores configurations dispatched by the dongle.
func (c *Client) QuickLogin(username, password string) (*LoginResp, error) {
	if err := c.UpdateSession(); err != nil {
		return nil, fmt.Errorf("could not renew session: %w", err)
	}

	login, err := c.GetLoginState()
	if err != nil {
		return nil, err
	}

	var opts RequestOptions = nil
	sess, err := c.GetSessionTokenInfo()
	if err != nil {
		return nil, fmt.Errorf("could not fetch session token: %w", err)
	}
	log.Printf("using salted password (token: %s)\n", sess.Token)
	encodedPassword := crypto.EncodeSaltedPassword(username, password, sess.Token)

	return c.doLogin(username, encodedPassword, login.PasswordType, opts)
}

// SlowLogin performs the login routine same as the browser.
func (c *Client) Login(username, password string) (*LoginResp, error) {
	if err := c.UpdateSession(); err != nil {
		return nil, fmt.Errorf("could not renew session: %w", err)
	}

	login, err := c.GetLoginState()
	if err != nil {
		return nil, err
	}

	if login.UseScarmLogin() {
		log.Println("unsupported login type: SCARM, Try to use quick login.")
		return c.QuickLogin(username, password)
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

	encodedPassword, err := c.encodePassword(username, password, login.IsPasswordSalted())
	if err != nil {
		return nil, err
	}

	return c.doLogin(username, encodedPassword, login.PasswordType, opts)
}

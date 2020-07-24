package client

import (
	"crypto/rsa"
	"encoding/xml"
	"math/big"
	"strconv"
)

type SessionTokenInfoResp struct {
	XMLName xml.Name `xml:"response"`
	Session string   `xml:"SesInfo"`
	Token   string   `xml:"TokInfo"`
}

func (c *Client) GetSessionTokenInfo() (*SessionTokenInfoResp, error) {
	session := SessionTokenInfoResp{}
	err := c.API("/webserver/SesTokInfo", nil, &session, nil)
	if err != nil {
		return nil, err
	}
	c.Token = session.Token
	return &session, nil
}

type PubKeyResp struct {
	XMLName xml.Name `xml:"response"`
	N       string   `xml:"encpubkeyn"`
	E       string   `xml:"encpubkeye"`
}

func (c *Client) GetPublicKey() (*rsa.PublicKey, error) {
	resp := PubKeyResp{}
	pubKey := &rsa.PublicKey{}
	err := c.API("/webserver/publickey", nil, &resp, nil)
	if err == nil {
		pubKey.E, _ = strconv.Atoi(resp.E)
		pubKey.N = &big.Int{}
		pubKey.N.SetString(resp.N, 16)
	}
	return pubKey, err
}

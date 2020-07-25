package client

import (
	"encoding/xml"
	"time"
)

type SMSCountResp struct {
	BaseResp

	InboxUnread  string `xml:"LocalUnread"`
	InboxTotal   string `xml:"LocalInbox"`
	OutboxCount  string `xml:"LocalOutbox"`
	DraftCount   string `xml:"LocalDraft"`
	DeletedCount string `xml:"LocalDeleted"`

	// Usually "500".
	// The max number of SMS can be stored in the device?
	LocalMaxSMSCount string `xml:"LocalMax"`

	// SMS Stored in the Sim card
	SimUnread string `xml:"SimUnread"`
	SimInbox  string `xml:"SimInbox"`
	SimOutbox string `xml:"SimOutbox"`
	SimDraft  string `xml:"SimDraft"`
	SimMax    string `xml:"SimMax"`
	SimUsed   string `xml:"SimUsed"`

	NewMsg string `xml:"NewMsg"`
}

// SMSCount: Count the number of SMS stored.
func (c *Client) SMSCount() (*SMSCountResp, error) {
	resp := &SMSCountResp{}
	err := c.API("/sms/sms-count", nil, resp, nil)
	return resp, err
}

type SMSMessage struct {
	XMLName xml.Name `xml:"Message"`

	// 0: Unread; 1: Read; 2: Draft; 3: Sent; 4: Warning/Error(?)
	Status int `xml:"Smstat"`

	MessageID string `xml:"Index"`

	// Can be multiple phone numbers, separated by ";"
	Phone   string `xml:"Phone"`
	Content string `xml:"Content"`
	Date    string `xml:"Date"`

	// Unknown property
	Sca string `xml:"Sca"`

	// Unknown value
	SaveType string `xml:"SaveType"`

	// Unknown Value
	Priority string `xml:"Priority"`

	// 7: Success, 8: Failed; other: skip
	SmsType int `xml:"SmsType"`
}

type SMSListResp struct {
	BaseResp

	Count    string       `xml:"Count"`
	Messages []SMSMessage `xml:"Messages>Message"`
}

type SMSListPayload struct {
	XMLName xml.Name `xml:"request"`

	// Page starts from 1
	Page int `xml:"PageIndex"`

	// default value in web is 20
	MessagesPerPage int `xml:"ReadCount"`

	// BoxType can be one of the following value:
	// LOCAL_INBOX = 1;
	// LOCAL_SENT  = 2;
	// LOCAL_DRAFT = 3;
	// LOCAL_TRASH = 4;
	// SIM_INBOX   = 5;
	// SIM_SENT    = 6;
	// SIM_DRAFT   = 7;
	// MIX_INBOX   = 8;
	// MIX_SENT    = 9;
	// MIX_DRAFT   = 10;
	BoxType         int `xml:"BoxType"`
	SortType        int `xml:"SortType"`
	Ascending       int `xml:"Ascending"`
	UnreadPreferred int `xml:"UnreadPreferred"`
}

// SMSList returns a list of SMS messages.
// page starts from 1.
// messagesPerPage is 20 in the web ui.
func (c *Client) SMSList(page, messagesPerPage int) (*SMSListResp, error) {
	body := &SMSListPayload{
		Page:            page,
		MessagesPerPage: messagesPerPage,
		BoxType:         1,
		SortType:        0,
		Ascending:       0,
		UnreadPreferred: 0,
	}

	return c.SMSListAPI(body)
}

func (c *Client) SMSListAPI(body *SMSListPayload) (*SMSListResp, error) {
	resp := &SMSListResp{}
	err := c.API("/sms/sms-list", body, resp, nil)
	return resp, err
}

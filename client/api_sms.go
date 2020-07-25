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

type SMSSendPayload struct {
	XMLName xml.Name `xml:"request"`

	// Set to -1; or a draft SMS ID
	ID int64 `xml:"Index"`

	// Phone numbers to sent
	Phones []string `xml:"Phones>Phone"`

	// empty string
	SCA string `xml:"Sca"`

	// Text message
	Content string `xml:"Content"`
	// Message size
	Length int `xml:"Length"`

	// Set this field to "1"
	SMSType int `xml:"Reserved"`

	// Format: "yyyy-MM-dd hh:mm:ss"
	// Use "SetDate" to help you.
	Date string `xml:"Date"`
}

func (s *SMSSendPayload) SetDate(date time.Time) {
	s.Date = date.Format("2006-01-02 15:04:05")
}

func (s *SMSSendPayload) SetDateToNow() {
	s.SetDate(time.Now())
}

type SMSSendResp struct {
	BaseResp
}

// SendSMS sends a text message to a specified phone number.
func (c *Client) SendSMS(phone, message string) (*SMSSendResp, error) {
	body := &SMSSendPayload{
		ID:      -1,
		Phones:  []string{phone},
		SCA:     "",
		Content: message,
		Length:  len(message),
		SMSType: 1,
	}
	body.SetDateToNow()
	return c.SendSMSAPI(body)
}

func (c *Client) SendSMSAPI(body *SMSSendPayload) (*SMSSendResp, error) {
	resp := &SMSSendResp{}
	err := c.API("/sms/send-sms", body, resp, nil)
	return resp, err
}

type SMSSendStatusResp struct {
	BaseResp

	// Phone number not delivered yet
	Phone string `xml:"Phone"`

	// Phone number delivered successfully
	SuccessPhone string `xml:"SucPhone"`

	// Phone number failed to deliver
	FailedPhone string `xml:"FailPhone"`

	// The number of message sent from last request
	TotalCount uint `xml:"TotalCount"`

	// Current index (start from 1)
	CurIndex uint `xml:"CurIndex"`
}

func (c *Client) GetSendStatus() (*SMSSendStatusResp, error) {
	resp := &SMSSendStatusResp{}
	err := c.API("/sms/send-status", nil, resp, nil)
	return resp, err
}

package client

import (
	"encoding/xml"
	"time"
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	SMS_UNREAD_STATUS = 0
	SMS_READ_STATUS   = 1
	SMS_DRAFT_STATUS  = 2
	SMS_SENT_STATUS   = 3
	SMS_ERROR_STATUS  = 4
)

// SMS box type
//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	SMS_BOX_TYPE_LOCAL_INBOX = 1
	SMS_BOX_TYPE_LOCAL_SENT  = 2
	SMS_BOX_TYPE_LOCAL_DRAFT = 3
	SMS_BOX_TYPE_LOCAL_TRASH = 4
	SMS_BOX_TYPE_SIM_INBOX   = 5
	SMS_BOX_TYPE_SIM_SENT    = 6
	SMS_BOX_TYPE_SIM_DRAFT   = 7
	SMS_BOX_TYPE_MIX_INBOX   = 8
	SMS_BOX_TYPE_MIX_SENT    = 9
	SMS_BOX_TYPE_MIX_DRAFT   = 10
)

type SMSCountResp struct {
	BaseResp

	InboxUnread  int `xml:"LocalUnread"`
	InboxTotal   int `xml:"LocalInbox"`
	OutboxCount  int `xml:"LocalOutbox"`
	DraftCount   int `xml:"LocalDraft"`
	DeletedCount int `xml:"LocalDeleted"`

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
	XMLName   xml.Name `xml:"Message"`
	Status    int      `xml:"Smstat"`
	MessageID string   `xml:"Index"`

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
		BoxType:         SMS_BOX_TYPE_LOCAL_INBOX,
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

type SMSSetReadPayload struct {
	XMLName xml.Name `xml:"request"`
	// Set to -1; or a draft SMS ID
	ID int64 `xml:"Index"`
}

func (c *Client) SetRead(Index int64) {
	body := &SMSSetReadPayload{
		ID: Index,
	}
	c.SetReadAPI(body)
}
func (c *Client) SetReadAPI(body *SMSSetReadPayload) {
	_ = c.API("/sms/set-read", body, nil, nil)
}

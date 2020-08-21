package client

import (
	"gopkg.in/guregu/null.v4/zero"
)

type NetworkType int

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	TYPE_NOSERVICE       = 0
	TYPE_GSM             = 1
	TYPE_GPRS            = 2
	TYPE_EDGE            = 3
	TYPE_WCDMA           = 4
	TYPE_HSDPA           = 5
	TYPE_HSUPA           = 6
	TYPE_HSPA            = 7
	TYPE_TDSCDMA         = 8
	TYPE_HSPA_PLUS       = 9
	TYPE_EVDO_REV_0      = 10
	TYPE_EVDO_REV_A      = 11
	TYPE_EVDO_REV_B      = 12
	TYPE_1xRTT           = 13
	TYPE_UMB             = 14
	TYPE_1xEVDV          = 15
	TYPE_3xRTT           = 16
	TYPE_HSPA_PLUS_64QAM = 17
	TYPE_HSPA_PLUS_MIMO  = 18
	TYPE_LTE             = 19
)

type NetworkTypeEx int

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	TYPE_EX_NOSERVICE          = 0
	TYPE_EX_GSM                = 1
	TYPE_EX_GPRS               = 2
	TYPE_EX_EDGE               = 3
	TYPE_EX_IS95A              = 21
	TYPE_EX_IS95B              = 22
	TYPE_EX_CDMA_1x            = 23
	TYPE_EX_EVDO_REV_0         = 24
	TYPE_EX_EVDO_REV_A         = 25
	TYPE_EX_EVDO_REV_B         = 26
	TYPE_EX_HYBRID_CDMA_1x     = 27
	TYPE_EX_HYBRID_EVDO_REV_0  = 28
	TYPE_EX_HYBRID_EVDO_REV_A  = 29
	TYPE_EX_HYBRID_EVDO_REV_B  = 30
	TYPE_EX_EHRPD_REL_0        = 31
	TYPE_EX_EHRPD_REL_A        = 32
	TYPE_EX_EHRPD_REL_B        = 33
	TYPE_EX_HYBRID_EHRPD_REL_0 = 34
	TYPE_EX_HYBRID_EHRPD_REL_A = 35
	TYPE_EX_HYBRID_EHRPD_REL_B = 36
	TYPE_EX_WCDMA              = 41
	TYPE_EX_HSDPA              = 42
	TYPE_EX_HSUPA              = 43
	TYPE_EX_HSPA               = 44
	TYPE_EX_HSPA_PLUS          = 45
	TYPE_EX_DC_HSPA_PLUS       = 46
	TYPE_EX_TD_SCDMA           = 61
	TYPE_EX_TD_HSDPA           = 62
	TYPE_EX_TD_HSUPA           = 63
	TYPE_EX_TD_HSPA            = 64
	TYPE_EX_TD_HSPA_PLUS       = 65
	TYPE_EX_802_16E            = 81
	TYPE_EX_LTE                = 101
)

type MonitorStatusResp struct {
	BaseResp

	ConnectionStatus     int         `xml:"ConnectionStatus"`
	WiFiConnectionStatus zero.Int    `xml:"WifiConnectionStatus"`
	SignalStrength       zero.Int    `xml:"SignalStrength"`
	SignalIcon           zero.Int    `xml:"SignalIcon"`
	CurrentNetworkType   NetworkType `xml:"CurrentNetworkType"`
	CurrentServiceDomain int         `xml:"CurrentServiceDomain"`
	RoamingStatus        int         `xml:"RoamingStatus"`

	// Battery settings
	BatteryStatus  string   `xml:"BatteryStatus"`
	BatteryLevel   zero.Int `xml:"BatteryLevel"`
	BatteryPercent zero.Int `xml:"BatteryPercent"`

	SIMLockStatus int `xml:"simlockStatus"`

	// Network settings
	PrimaryDns       string `xml:"PrimaryDns"`
	SecondaryDns     string `xml:"SecondaryDns"`
	PrimaryIPv6Dns   string `xml:"PrimaryIPv6Dns"`
	SecondaryIPv6Dns string `xml:"SecondaryIPv6Dns"`

	// WiFi Related
	CurrentWiFiUser      string        `xml:"CurrentWifiUser"`
	TotalWiFiUser        string        `xml:"TotalWifiUser"`
	CurrentTotalWiFiUser uint          `xml:"currenttotalwifiuser"`
	ServiceStatus        string        `xml:"ServiceStatus"`
	SimStatus            string        `xml:"SimStatus"`
	WiFiStatus           string        `xml:"WifiStatus"`
	CurrentNetworkTypeEx NetworkTypeEx `xml:"CurrentNetworkTypeEx"`
	WanPolicy            string        `xml:"WanPolicy"`
	MaxSignal            string        `xml:"maxsignal"`
	WiFiIndoorOnly       int           `xml:"wifiindooronly"`
	WiFiFrequency        int           `xml:"wififrequence"`

	// Can be one of: "mobile-wifi", "cpe", "hilink".
	// Maybe other values as well?
	DeviceType string `xml:"classify"`

	// Network?
	AirplaneMode int `xml:"flymode"`
	CellRoaming  int `xml:"cellroam"`

	// Optional?
	VoiceBusy zero.Int `xml:"voice_busy"`
	UsbUP     zero.Int `xml:"usbup"`
}

type NetworkPLMNResp struct {
	BaseResp

	State     string `xml:"State"`
	FullName  string `xml:"FullName"`
	ShortName string `xml:"ShortName"`
	Numeric   int    `xml:"Numeric"`
	Rat       int    `xml:"Rat"`
	Spn       string `xml:"Spn"`
}

func (c *Client) GetNetworkPLMN() (*NetworkPLMNResp, error) {
	resp := &NetworkPLMNResp{}
	err := c.API("/net/current-plmn", nil, resp, nil)
	return resp, err
}

func (c *Client) GetDeviceStatus() (*MonitorStatusResp, error) {
	resp := &MonitorStatusResp{}
	err := c.API("/monitoring/status", nil, resp, nil)
	return resp, err
}

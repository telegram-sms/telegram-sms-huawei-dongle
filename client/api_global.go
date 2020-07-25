package client

type ModuleSwitchResp struct {
	BaseResp

	USSDEnabled          string `xml:"ussd_enabled"`
	BBOUEnabled          string `xml:"bbou_enabled"`
	SMSEnabled           string `xml:"sms_enabled"`
	SDCardEnabled        string `xml:"sdcard_enabled"`
	WiFiEnabled          string `xml:"wifi_enabled"`
	StatisticEnabled     string `xml:"statistic_enabled"`
	HelpEnabled          string `xml:"help_enabled"`
	STKEnabled           string `xml:"stk_enabled"`
	PBEnabled            string `xml:"pb_enabled"`
	DLNAEnabled          string `xml:"dlna_enabled"`
	OTAEnabled           string `xml:"ota_enabled"`
	WiFiOffloadEnabled   string `xml:"wifioffload_enabled"`
	CradleEnabled        string `xml:"cradle_enabled"`
	MultiSSIDEnable      string `xml:"multssid_enable"`
	IPv6Enabled          string `xml:"ipv6_enabled"`
	MonthlyVolumeEnabled string `xml:"monthly_volume_enabled"`
	PowerSaveEnabled     string `xml:"powersave_enabled"`
	SNTPEnabled          string `xml:"sntp_enabled"`
	EncryptEnabled       string `xml:"encrypt_enabled"`
	DataSwitchEnabled    string `xml:"dataswitch_enabled"`
	PowerOffEnabled      string `xml:"poweroff_enabled"`
	EcoModeEnabled       string `xml:"ecomode_enabled"`
	ZoneTimeEnabled      string `xml:"zonetime_enabled"`
	LocalUpdateEnabled   string `xml:"localupdate_enabled"`
	CBSEnabled           string `xml:"cbs_enabled"`
	QRCodeEnabled        string `xml:"qrcode_enabled"`
	ChargerEnabled       string `xml:"charger_enbaled"`
	APNRetryEnabled      string `xml:"apn_retry_enabled"`
	GDPREnabled          string `xml:"gdpr_enabled"`
}

func (r *ModuleSwitchResp) IsEncryptionEnabled() bool {
	return r.EncryptEnabled == "1"
}

func (c *Client) GetModuleSwitches() (*ModuleSwitchResp, error) {
	resp := &ModuleSwitchResp{}
	err := c.API("/global/module-switch", nil, resp, nil)
	return resp, err
}

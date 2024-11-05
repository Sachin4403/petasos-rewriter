package main

type UpdateResourceRequest struct {
	IpAddress               string `json:"ipAddress"`
	CertificateProviderType string `json:"certificateProviderType"`
	CertificateExpiryDate   string `json:"certificateExpiryDate"`
	LastRebootReason        string `json:"lastRebootReason"`
	WanInterfaceUsed        string `json:"wanInterfaceUsed"`
	LastReconnectReason     string `json:"lastReconnectReason"`
	ManagementProtocol      string `json:"managementProtocol"`
	LastBootTime            int64  `json:"lastBootTime"`
	FirmwareVersion         string `json:"firmwareVersion"`
}

type WebPAConveyHeaderData struct {
	WebpaProtocol            string `json:"webpa-protocol"`
	WebpaInterfaceUsed       string `json:"webpa-interface-used"`
	HwLastRebootReason       string `json:"hw-last-reboot-reason"`
	WebpaLastReconnectReason string `json:"webpa-last-reconnect-reason"`
	BootTime                 int64  `json:"boot-time"`
	FwName                   string `json:"fw-name"`
}

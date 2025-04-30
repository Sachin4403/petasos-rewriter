package main

type UpdateResourceRequest struct {
	IpAddress                string `json:"ipAddress"`
	CertificateProviderType  string `json:"certificateProviderType"`
	CertificateExpiryDate    string `json:"certificateExpiryDate"`
	LastRebootReason         string `json:"lastRebootReason,omitempty"`
	WanInterfaceUsed         string `json:"wanInterfaceUsed,omitempty"`
	LastReconnectReason      string `json:"lastReconnectReason,omitempty"`
	ManagementProtocol       string `json:"managementProtocol,omitempty"`
	FirmwareVersion          string `json:"firmwareVersion,omitempty"`
	Ipv4AddressHGWWAN        string `json:"ipv4AddressHGWWAN,omitempty"`
	IsPetasosRewriterRequest bool   `json:"isPetasosRewriterRequest,omitempty"`
}

type WebPAConveyHeaderData struct {
	WebpaProtocol            string `json:"webpa-protocol"`
	WebpaInterfaceLabel      string `json:"webpa-interface-label"`
	HwLastRebootReason       string `json:"hw-last-reboot-reason"`
	WebpaLastReconnectReason string `json:"webpa-last-reconnect-reason"`
	BootTime                 int64  `json:"boot-time"`
	FwName                   string `json:"fw-name"`
	Ipv4AddressHGWWAN        string `json:"wan-ipv4-address"`
}

package main

type UpdateResourceRequest struct {
	IpAddress                string `json:"ipAddress"`
	CertificateProviderType  string `json:"certificateProviderType"`
	CertificateExpiryDate    string `json:"certificateExpiryDate"`
	HwLastRebootReason       string `json:"hw-last-reboot-reason"`
	WebpaInterfaceUsed       string `json:"webpa-interface-used"`
	WebpaLastReconnectReason string `json:"webpa-last-reconnect-reason"`
	BootTime                 string `json:"bootTime"`
	WebpaProtocol            string `json:"webpa-protocol"`
}

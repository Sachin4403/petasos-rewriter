package main

import "time"

type UpdateResourceRequestForRI struct {
	IpAddress               string    `json:"ipAddress"`
	CertificateProviderType string    `json:"certificateProviderType"`
	CertificateExpiryDate   string    `json:"certificateExpiryDate"`
	LastRebootReason        string    `json:"lastRebootReason"`
	WanInterfaceUsed        string    `json:"wanInterfaceUsed"`
	LastReconnectReason     string    `json:"lastReconnectReason"`
	ManagementProtocol      string    `json:"managementProtocol"`
	LastBootTime            time.Time `json:"lastBootTime"`
	FirmwareVersion         string    `json:"firmwareVersion"`
}

type WebPAHeaderData struct {
	WebpaProtocol            string `json:"webpa-protocol"`
	WebpaInterfaceUsed       string `json:"webpa-interface-used"`
	HwLastRebootReason       string `json:"hw-last-reboot-reason"`
	WebpaLastReconnectReason string `json:"webpa-last-reconnect-reason"`
	BootTime                 int64  `json:"boot-time"`
	FwName                   string `json:"fw-name"`
}

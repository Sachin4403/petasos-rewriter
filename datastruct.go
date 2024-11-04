package main

import "time"

type UpdateResourceRequest struct {
	IpAddress                string    `json:"ipAddress"`
	CertificateProviderType  string    `json:"certificateProviderType"`
	CertificateExpiryDate    string    `json:"certificateExpiryDate"`
	HwLastRebootReason       string    `json:"lastRebootReason"`
	WebpaInterfaceUsed       string    `json:"wanInterfaceUsed"`
	WebpaLastReconnectReason string    `json:"lastReconnectReason"`
	WebpaProtocol            string    `json:"managementProtocol"`
	LastBootTime             time.Time `json:"lastBootTime"`
	FirmwareVersion          string    `json:"firmwareVersion"`
}
type WebPAHeaderData struct {
	WebpaProtocol            string `json:"webpa-protocol"`
	WebpaInterfaceUsed       string `json:"webpa-interface-used"`
	HwLastRebootReason       string `json:"hw-last-reboot-reason"`
	WebpaLastReconnectReason string `json:"webpa-last-reconnect-reason"`
	BootTime                 int64  `json:"boot-time"`
	FwName                   string `json:"fw-name"`
}

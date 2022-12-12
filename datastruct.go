package main

type UpdateResourceRequest struct {
	Cnmac       string `json:"cnmac"`
	Environment string `json:"environment"`
	Mac         string `json:"mac"`
	RemoteIp    string `json:"remoteIp"`
	TenantId    string `json:"x-tenant-id"`
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestReplaceTalariaInternalName(t *testing.T) {

	testData := []struct {
		host     string
		old      string
		new      string
		expected string
		err      error
	}{
		{"xmidt-talaria-1", "xmidt-talaria-", "talaria", "talaria1", nil},
		{"xmidt-talaria-2", "xmidt-talaria-", "talaria", "talaria2", nil},
		{
			host:     "xmidt-talaria3",
			old:      "xmidt-talaria",
			new:      "talaria",
			expected: "talaria3",
		},
		{
			host:     "xmidt-talaria4",
			old:      "xmidt-talaria",
			new:      "talaria",
			expected: "talaria4",
		},
		{"xmidt-talaria4", "xmidt-talaria-", "talaria", "talaria4", ErrNoMatchFound},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert      = assert.New(t)
				actual, err = replaceTalariaInternalName(record.host, record.old, record.new)
			)
			if err != nil {
				assert.Equal(record.err, err)
			} else {
				assert.Equal(record.expected, actual)
			}
		})
	}

}

func TestBuildExternalURL(t *testing.T) {
	testData := []struct {
		arg1     string
		arg2     string
		expected string
	}{
		{"", "", "."},
		{"talaria", "Test.com", "talaria.Test.com"},
		{"talaria2", "dev.rdk.yo-digital.com", "talaria2.dev.rdk.yo-digital.com"},
		{"talaria3", "xyz.com", "talaria3.xyz.com"},
	}
	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert = assert.New(t)
				actual = buildExternalURL(record.arg1, record.arg2)
			)
			assert.Equal(record.expected, actual)

		})
	}

}

func TestForwarder(t *testing.T) {
	testData := []struct {
		deviceName string
	}{
		{"mac:B827EBB25F81"},
		{"mac:B827EBB25F82"},
		{"mac:B827EBB25F83"},
		{"mac:B827EBB25F84"},
		{"mac:B827EBB25F85"},
		{"mac:B827EBB25F86"},
		{"mac:B827EBB25F87"},
	}
	e := echo.New()

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert  = assert.New(t)
				handler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
					assert.Equal(record.deviceName, request.Header.Get("X-Webpa-Device-Name"))
					response.Header().Set("Content-Type", "text/html; charset=utf-8")
					response.Header().Set("Location", "http://xmidt-talaria:6200/api/v2/device")
					response.Header().Set("Date", time.Now().String())
					response.Header().Set("X-Petasos-Build", "Test")
					response.Header().Set("X-Petasos-Flavor", "Test")
					response.Header().Set("X-Petasos-Region", "Test")
					response.Header().Set("X-Petasos-Server", "Test")
					response.Header().Set("X-Webpa-Device-Name", record.deviceName)
					body := "<a href=\"http://xmidt-talaria:6200/api/v2/device\">Temporary Redirect</a>.\n"
					response.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
					response.Write([]byte(body))
				})
			)

			server := httptest.NewServer(handler)
			defer server.Close()
			petasosURL, _ = url.Parse(server.URL)
			r := httptest.NewRequest("", "/v2/api/device", nil)
			r.Header.Set("X-Webpa-Device-Name", record.deviceName)
			r.Header.Set("X-Forwarded-Proto", "ws")
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)
			err := forwarder(c, client)
			assert.Nil(err)

		})
	}
}
func TestUpdateResourceIpAddressAndCertificateInfo(t *testing.T) {
	testsData := []struct {
		realIP                   string
		certificateProvider      string
		expiryDate               string
		deviceCN                 string
		hwLastRebootReason       string
		webpaInterfaceUsed       string
		webpaLastReconnectReason string
		webpaProtocol            string
		expectedRequestBody      string
		expectedStatus           int
		xWebPA                   string
	}{
		{
			realIP:                   "127.0.0.1",
			certificateProvider:      "DTSECURITY",
			expiryDate:               "Sep 19 23:59:59 2031 GMT",
			deviceCN:                 "TestCPE",
			hwLastRebootReason:       "unknown",
			webpaInterfaceUsed:       "erouter0",
			webpaLastReconnectReason: "SSL_Socket_Close",
			webpaProtocol:            "PARODUS-2.0-61b1a7a",
			xWebPA:                   "eyJody1tb2RlbCI6IlwiRkdBMjIzM1wiIiwiaHctc2VyaWFsLW51bWJlciI6IjIyMzNBRENNTCIsImh3LW1hbnVmYWN0dXJlciI6IlwiVGVjaG5pY29sb3JcIiIsImZ3LW5hbWUiOiIwMDUuMDMzLjAwMSIsImJvb3QtdGltZSI6MTcyNTAwMDYwOCwid2VicGEtcHJvdG9jb2wiOiJQQVJPRFVTLTIuMC02MWIxYTdhIiwid2VicGEtaW50ZXJmYWNlLXVzZWQiOiJlcm91dGVyMCIsImh3LWxhc3QtcmVib290LXJlYXNvbiI6InVua25vd24iLCJ3ZWJwYS1sYXN0LXJlY29ubmVjdC1yZWFzb24iOiJTU0xfU29ja2V0X0Nsb3NlIn0=",
			expectedRequestBody:      `{"ipAddress":"127.0.0.1","certificateProviderType":"DTSECURITY","certificateExpiryDate":"Sep 19 23:59:59 2031 GMT","hw-last-reboot-reason":"unknown","webpa-interface-used":"erouter0","webpa-last-reconnect-reason":"SSL_Socket_Close","webpa-protocol":"PARODUS-2.0-61b1a7a"}`,
			expectedStatus:           http.StatusOK,
		},
	}

	for i, tt := range testsData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert := assert.New(t)

			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Log the received headers for debugging
				t.Logf("Received headers: %+v", r.Header)

				assert.Equal(http.MethodPut, r.Method)

				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(err)
				assert.JSONEq(tt.expectedRequestBody, string(requestBody))
				w.WriteHeader(tt.expectedStatus)
			}))
			defer mockServer.Close()

			testReq, err := http.NewRequest(http.MethodPut, "/", nil)
			assert.NoError(err)
			testReq.Header.Set(realIpHeader, tt.realIP)
			testReq.Header.Set(certificateProviderHeader, tt.certificateProvider)
			testReq.Header.Set(expiryDateHeader, tt.expiryDate)
			testReq.Header.Set(deviceCNHeader, tt.deviceCN)
			testReq.Header.Set(lastRebootReason, tt.hwLastRebootReason)
			testReq.Header.Set(webpaInterfaceUsed, tt.webpaInterfaceUsed)
			testReq.Header.Set(webpaLastReconnectReason, tt.webpaLastReconnectReason)
			testReq.Header.Set(webpaProtocol, tt.webpaProtocol)
			testReq.Header.Set("ENVIRONMENT", "test")
			testReq.Header.Set("X-TENANT-ID", "12345")
			testReq.Header.Set(xWebPA, tt.xWebPA)

			client := &http.Client{
				Transport: &http.Transport{
					Proxy: func(req *http.Request) (*url.URL, error) {
						return url.Parse(mockServer.URL)
					},
				},
			}

			resourceURL, err := url.Parse(mockServer.URL + "/v1/resource/macAddress")
			assert.NoError(err)

			err = updateResourceIpAddressAndCertificateInfo(testReq, client, resourceURL)

			if tt.expectedStatus != http.StatusOK {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(http.StatusOK, tt.expectedStatus, "Expected status OK")
			} else {
				assert.NotEqual(http.StatusOK, tt.expectedStatus, "Expected status to be not OK")
			}
		})
	}
}

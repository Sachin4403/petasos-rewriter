package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
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
			assert := assert.New(t)
			mux := http.NewServeMux()
			mux.HandleFunc("/v2/api/device", func(response http.ResponseWriter, request *http.Request) {
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

			mux.HandleFunc("/hgw/rdk-xmidt-executor/v1/listener/xmidt/device", func(response http.ResponseWriter, request *http.Request) {
				decoder := json.NewDecoder(request.Body)
				var updateResourceRequest UpdateResourceRequest
				decoder.Decode(&updateResourceRequest)
				assert.Equal(record.deviceName, updateResourceRequest.Cnmac)
				assert.Equal(record.deviceName, updateResourceRequest.Mac)
				assert.Equal("dev", updateResourceRequest.Environment)
				assert.Equal("11.111.111.11", updateResourceRequest.RemoteIp)
				assert.Equal("dtdl", updateResourceRequest.TenantId)
			})

			server := httptest.NewServer(mux)
			defer server.Close()
			petasosURL, _ = url.Parse(server.URL)
			r := httptest.NewRequest("", "/v2/api/device", nil)
			r.Header.Set("X-Webpa-Device-Name", record.deviceName)
			r.Header.Set("X-Forwarded-Proto", "ws")
			r.Header.Set("x-device-cn", record.deviceName)
			r.Header.Set("remoteIp", "11.111.111.11")
			r.Header.Set("x-tenant-id", "dtdl")
			r.Header.Set("environment", "dev")
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

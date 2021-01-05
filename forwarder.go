package main

import (
	"fmt"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

var (
	ErrNoMatchFound = fmt.Errorf("No match found")
	spanIdHeader = "X-B3-SpanId"
	traceIdHeader = "X-B3-TraceId"
)

// forwarder forwads requests to real petasos instance and does
// apropriate replacements.
func forwarder(c echo.Context) error {
	sp := jaegertracing.CreateChildSpan(c, "got request in forwarder.go")
	defer  sp.Finish()
	jaegarCtx := sp.Context().(jaeger.SpanContext)
	traceId := jaegarCtx.TraceID().String()
	spanId := jaegarCtx.SpanID().String()
	debugLog := log.Debug().Str(spanIdHeader,spanId).Str(traceIdHeader,traceId)
	infoLog := log.Info().Str(spanIdHeader,spanId).Str(traceIdHeader,traceId)
	debugLog.Msg("##############################")
	debugLog.Msg("###### Request Start #########")
	debugLog.Msg("##############################")

	// prepare request for forwarding
	req := c.Request()
	//deviceId := req.Header.Get("X-Webpa-Device-Name")


	//sp.LogKV(
	//	"device-id",deviceId,
	//)
	//sp.SetBaggageItem("device-name",deviceId)
	//sp.SetTag("application","petasos-rewriter")


	// store scheme of original request
	originalRequestScheme := req.URL.Scheme
	if originalRequestScheme == "" {
		originalRequestScheme = req.Header.Get("X-Forwarded-Proto")
	}
	debugLog.Msgf("Trace_id %s and spanId %s",traceId,spanId)
	debugLog.Msgf("originalScheme [%s]", originalRequestScheme)

	// Change protocols fro`m ws(s) => http(s).
	// Parodus makes requests to `ws` but complains
	// when getting a redirect containing `ws`.
	switch originalRequestScheme {
	case "ws":
		debugLog.Msgf("Replacing original scheme [%s] with [%s] in output", originalRequestScheme, "http")
		originalRequestScheme = "http"
	case "wss":
		debugLog.Msgf("Replacing original scheme [%s] with [%s] in output", originalRequestScheme, "https")
		originalRequestScheme = "https"
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	debugLog.Msg("Dumping original request to petasos-rewriter")
	debugLog.Msgf("%s", dump)
	debugLog.Msg("") // br
	debugLog.Msg("") // br

	// Prepare forwarding to petasos
	req.URL = &url.URL{
		Scheme: petasosURL.Scheme,
		Host:   petasosURL.Host,
		Path:   req.URL.Path,
	}
	req.RequestURI = ""

	// Forward to real petasos
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	opentracing.GlobalTracer().Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	dump, err = httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	debugLog.Msg("Dumping request to real petasos")
	debugLog.Msgf("%s", dump)
	debugLog.Msg("") // br
	debugLog.Msg("") // br


	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	debugLog.Msg("Dumping response from real petasos")
	debugLog.Msgf("%s", dump)
	debugLog.Msg("") // br
	debugLog.Msg("") // br

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// just printing the all response headers which we got from actual petasos
	for k, v := range resp.Header {
		var header string
		for _, s := range v {
			if header != "" {
				header = header + ","
			}
			header = header + s
		}
		header = strings.TrimRight(header, ",")
		c.Response().Header().Set(k, header)

		debugLog.Msgf("k: %s, v: %s\n", k, v)
	}



	if resp.StatusCode != http.StatusTemporaryRedirect {
		// Forward status code
		c.Response().Writer.WriteHeader(resp.StatusCode)
		c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		c.Response().Writer.Write(body)
		return nil
	}
	// Replace location header
	location := c.Response().Header().Get("Location")
	debugLog.Msgf("Location [%s]\n", location)

	locationUrl, err := url.Parse(location)
	if err != nil {
		return err
	}

	if *fixedScheme != "" {
		// TODO: use scheme from publicTalariaURL and make fixedScheme bool
		// locationUrl.Scheme = publicTalariaURL.Scheme
		locationUrl.Scheme = *fixedScheme
	} else {
		locationUrl.Scheme = originalRequestScheme
	}

	// Do replacement & build public talaria url
	externalTalariaName, err := replaceTalariaInternalName(
		locationUrl.Hostname(),
		*talariaInternalName,
		*talariaExternalName,
	)
	if err != nil {
		return err
	}
	publicTalariaURL := buildExternalURL(externalTalariaName, *talariaDomain)

	locationUrl.Host = publicTalariaURL
	infoLog.Msgf("redirecting from Location [%s] to Location [%s] for device name [%s] \n", location, locationUrl.String(),req.Header.Get("X-Webpa-Device-Name"))
	c.Response().Header().Set("Location", locationUrl.String())

	// Replace url in body
	var href = regexp.MustCompile(`"(.*)"`)
	body = href.ReplaceAll(body, []byte(`"`+locationUrl.String()+`"`))
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	c.Response().Header().Set(spanIdHeader, spanId)
	c.Response().Header().Set(traceIdHeader, traceId)

	// Forward status code
	c.Response().Writer.WriteHeader(resp.StatusCode)

	_, err = c.Response().Writer.Write(body)
	if err != nil {
		return err
	}


	return nil
}

// replaceTalariaInternalName replaces internal talaria name.
// Returns a ErrNoMatchFound when replacement is impossible.
func replaceTalariaInternalName(host, old, new string) (string, error) {
	index := strings.Index(host, old)
	if index == -1 {
		return "", ErrNoMatchFound
	}
	talariaExternal := strings.Replace(host, old, new, -1)

	// TODO: strip possible internal k8s namespace.
	// xmidt-talaria OK
	// talaria.xmidt Not OK

	return talariaExternal, nil
}

// buildExternalURL by concatenation new talaria name + given domain
func buildExternalURL(newTalariaName, domain string) string {
	var builder strings.Builder
	builder.WriteString(newTalariaName)
	builder.WriteString(".")
	builder.WriteString(domain)
	return builder.String()
}

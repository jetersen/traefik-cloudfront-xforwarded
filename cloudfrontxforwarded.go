package cloudfrontxforwarded

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	CloudFrontViewerAddressHeader  = "CloudFront-Viewer-Address"
	CloudFrontForwardedProtoHeader = "CloudFront-Forwarded-Proto"
	XForwardedForHeader            = "X-Forwarded-For"
	XForwardedProtoHeader          = "X-Forwarded-Proto"
	XForwardedPortHeader           = "X-Forwarded-Port"
	XRealIPHeader                  = "X-Real-IP"
	HttpsProtocol                  = "https"
)

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type CloudFrontXForwarded struct {
	next   http.Handler
	name   string
	config *Config
}

// Create a new plugin instance for CloudFrontXForwarded middleware.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &CloudFrontXForwarded{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

// CloudFrontXForwarded is a middleware that read CloudFront headers and sets the X-Forwarded-* headers accordingly.
func (cf *CloudFrontXForwarded) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Read the CloudFront headers and set the X-Forwarded-* headers.
	remoteAddressPort := req.Header.Get(CloudFrontViewerAddressHeader)
	if remoteAddressPort == "" {
		// If the CloudFront header is not set, just call the next handler.
		cf.next.ServeHTTP(rw, req)
		return
	}

	remoteAddr, remotePort, err := extractAddressPort(remoteAddressPort)
	if err != nil {
		// If the CloudFront header is not set, just call the next handler.
		cf.next.ServeHTTP(rw, req)
		return
	}

	req.Header.Set(XForwardedForHeader, remoteAddressPort)
	req.Header.Set(XRealIPHeader, remoteAddr)
	req.Header.Set(XForwardedPortHeader, remotePort)

	forwardedProto := req.Header.Get(CloudFrontForwardedProtoHeader)
	if forwardedProto != "" {
		req.Header.Set(XForwardedProtoHeader, forwardedProto)
	}

	cf.next.ServeHTTP(rw, req)
}

// extractAddressPort extracts the address and port from the CloudFront viewer address header and assumes format is "IP:port"
func extractAddressPort(addr string) (string, string, error) {
	colonIndex := strings.LastIndex(addr, ":")
	if colonIndex == -1 || colonIndex == len(addr)-1 {
		// return error as cloudfront-viewer-address should always have port
		return "", "", fmt.Errorf("cloudfront-viewer-address should always have port")
	}
	// Check if port can be parsed as integer
	port := addr[colonIndex+1:]
	if _, err := strconv.Atoi(port); err != nil {
		return "", "", fmt.Errorf("cloudfront-viewer-address port must be a number")
	}

	return addr[:colonIndex], port, nil
}

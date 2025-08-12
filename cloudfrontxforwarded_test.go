package cloudfrontxforwarded_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudfrontxforwarded "github.com/jetersen/traefik-cloudfront-xforwarded"
)

func TestCloudFrontXForwarded_ServeHTTP(t *testing.T) {
	ctx := context.Background()

	nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	handler, err := cloudfrontxforwarded.New(ctx, nextHandler, cloudfrontxforwarded.CreateConfig(), "cloudfront")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		url         string
		viewerAddr  string
		proto       string
		wantXFF     string
		wantXRealIP string
		wantXFProto string
		wantXFPort  string
	}{
		{
			name:        "ipv4 https",
			url:         "https://example.com",
			viewerAddr:  "1.2.3.4:12345",
			proto:       "https",
			wantXFF:     "1.2.3.4:12345",
			wantXRealIP: "1.2.3.4",
			wantXFProto: "https",
			wantXFPort:  "12345",
		},
		{
			name:        "ipv6 https",
			url:         "https://example.com",
			viewerAddr:  "2001:db8::1:54321",
			proto:       "https",
			wantXFF:     "2001:db8::1:54321",
			wantXRealIP: "2001:db8::1",
			wantXFProto: "https",
			wantXFPort:  "54321",
		},
		{
			name:        "missing viewer address header",
			url:         "https://example.com",
			viewerAddr:  "",
			proto:       "https",
			wantXFF:     "",
			wantXRealIP: "",
			wantXFProto: "",
			wantXFPort:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.viewerAddr != "" {
				req.Header.Set(cloudfrontxforwarded.CloudFrontViewerAddressHeader, tc.viewerAddr)
			}
			if tc.proto != "" {
				req.Header.Set(cloudfrontxforwarded.CloudFrontForwardedProtoHeader, tc.proto)
			}

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assertHeader(t, req, cloudfrontxforwarded.XForwardedForHeader, tc.wantXFF)
			assertHeader(t, req, cloudfrontxforwarded.XRealIPHeader, tc.wantXRealIP)
			assertHeader(t, req, cloudfrontxforwarded.XForwardedProtoHeader, tc.wantXFProto)
			assertHeader(t, req, cloudfrontxforwarded.XForwardedPortHeader, tc.wantXFPort)
		})
	}
}

func assertHeader(t *testing.T, req *http.Request, header, expectedValue string) {
	t.Helper()
	value := req.Header.Get(header)
	if value != expectedValue {
		t.Errorf("Expected %s header to be '%s', got '%s'", header, expectedValue, value)
	}
}

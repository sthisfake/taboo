package handler

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	allowedPrefixes = []string{"/vless"}

	stripHeaders = map[string]struct{}{
		"host":                {},
		"proxy-authenticate":  {},
		"proxy-authorization": {},
		"te":                  {},
		"trailer":             {},
		"transfer-encoding":   {},
		"forwarded":           {},
		"x-forwarded-host":    {},
		"x-forwarded-proto":   {},
		"x-forwarded-port":    {},
	}

	targetURL *url.URL
	proxy     *httputil.ReverseProxy
)

func init() {
	targetBase := strings.TrimSpace(os.Getenv("TARGET_DOMAIN"))
	if targetBase == "" {
		log.Fatal("TARGET_DOMAIN is not set")
	}

	var err error
	targetURL, err = url.Parse(targetBase)
	if err != nil {
		log.Fatalf("invalid TARGET_DOMAIN: %v", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		ForceAttemptHTTP2: false,

		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     0,

		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,

		ExpectContinueTimeout: 1 * time.Second,

		DisableCompression: true,

		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	proxy = &httputil.ReverseProxy{
		Transport: transport,

		Director: func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host

			req.Host = targetURL.Host

			originalHeaders := req.Header.Clone()

			req.Header = make(http.Header)
			copyHeaders(req.Header, originalHeaders)

			// preserve websocket upgrade headers
			if strings.EqualFold(originalHeaders.Get("Upgrade"), "websocket") {
				req.Header.Set("Connection", "Upgrade")
				req.Header.Set("Upgrade", "websocket")
			}

			clientIP := clientRealIP(req)
			if clientIP != "" {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		},

		FlushInterval: 10 * time.Millisecond,

		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error: %v", err)
			http.Error(rw, "Bad Gateway", http.StatusBadGateway)
		},
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}

	if !isAllowedPath(r.URL.Path) {
		http.NotFound(w, r)
		return
	}

	proxy.ServeHTTP(w, r)
}

func isAllowedPath(path string) bool {
	for _, p := range allowedPrefixes {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		lk := strings.ToLower(k)

		if _, blocked := stripHeaders[lk]; blocked {
			continue
		}

		if strings.HasPrefix(lk, "x-vercel-") {
			continue
		}

		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func clientRealIP(r *http.Request) string {
	ip := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return ""
}

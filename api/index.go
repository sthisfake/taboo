package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	stripHeaders = map[string]struct{}{
		"host":                {},
		"connection":          {},
		"keep-alive":          {},
		"proxy-authenticate":  {},
		"proxy-authorization": {},
		"te":                  {},
		"trailer":             {},
		"transfer-encoding":   {},
		"upgrade":             {},
		"forwarded":           {},
		"x-forwarded-host":    {},
		"x-forwarded-proto":   {},
		"x-forwarded-port":    {},
	}
	allowedPrefixes = []string{"/xhttp", "/vless"}
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}

	targetBase := strings.TrimSpace(os.Getenv("TARGET_DOMAIN"))
	if targetBase == "" {
		http.Error(w, "Misconfigured: TARGET_DOMAIN is not set", http.StatusInternalServerError)
		return
	}

	if !isAllowedPath(r.URL.Path) {
		http.NotFound(w, r)
		return
	}

	target, err := url.Parse(targetBase)
	if err != nil {
		http.Error(w, "Bad TARGET_DOMAIN", http.StatusBadGateway)
		return
	}
	targetURL := *target
	targetURL.Path = r.URL.Path
	targetURL.RawPath = r.URL.EscapedPath()
	targetURL.RawQuery = r.URL.RawQuery
	if _, err := url.Parse(targetURL.String()); err != nil {
		http.Error(w, "Bad TARGET_DOMAIN or request URL", http.StatusBadGateway)
		return
	}

	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = targetURL.Path
			req.URL.RawPath = targetURL.RawPath
			req.URL.RawQuery = targetURL.RawQuery
			req.Host = targetURL.Host

			// Build upstream headers from the original request while removing problematic ones.
			req.Header = make(http.Header)
			copyHeaders(req.Header, r.Header)

			clientIP := strings.TrimSpace(r.Header.Get("x-real-ip"))
			if clientIP == "" {
				clientIP = strings.TrimSpace(r.Header.Get("x-forwarded-for"))
			}
			if clientIP != "" {
				req.Header.Set("x-forwarded-for", clientIP)
			}
		},
		Transport: &http.Transport{
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		FlushInterval: -1, // flush immediately for streaming-style transports
		ErrorHandler: func(rw http.ResponseWriter, _ *http.Request, _ error) {
			http.Error(rw, "Bad Gateway: Tunnel Failed", http.StatusBadGateway)
		},
	}
	rp.ServeHTTP(w, r)
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

// reverse proxy copies response headers/body itself.

package httputil

import (
	"net/http"
	"net/url"
	"strings"
)

// AdminMutationAllowed rejects browser cross-site state-changing admin calls.
// Non-browser clients (no Origin/Referer) remain allowed for local curl/CLI use.
func AdminMutationAllowed(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	referer := strings.TrimSpace(r.Header.Get("Referer"))
	if origin == "" && referer == "" {
		return true
	}
	host := strings.TrimSpace(r.Host)
	if fwd := strings.TrimSpace(r.Header.Get("X-Forwarded-Host")); fwd != "" {
		host = fwd
	}
	if origin != "" {
		return hostMatchesURL(origin, host)
	}
	return hostMatchesURL(referer, host)
}

func hostMatchesURL(raw, host string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return false
	}
	return strings.EqualFold(u.Host, host)
}

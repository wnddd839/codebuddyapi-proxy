package httputil

import (
	"net/http"
	"net/url"
	"strings"
)

// AdminMutationAllowed 拒绝浏览器跨站发起的管理台写操作。
// 无 Origin/Referer 的非浏览器客户端（curl/CLI）仍允许，便于本地调试。
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

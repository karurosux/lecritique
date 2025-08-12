package utils

import (
	"net"
	"net/http"
	"strings"
)

type DeviceInfo struct {
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
	Platform  string `json:"platform"`
	Browser   string `json:"browser"`
}

func ExtractDeviceInfo(r *http.Request) *DeviceInfo {
	userAgent := r.Header.Get("User-Agent")
	ip := getClientIP(r)
	platform := extractPlatform(userAgent)
	browser := extractBrowser(userAgent)

	return &DeviceInfo{
		UserAgent: userAgent,
		IP:        ip,
		Platform:  platform,
		Browser:   browser,
	}
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func extractPlatform(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "windows") {
		return "Windows"
	}
	if strings.Contains(ua, "mac") || strings.Contains(ua, "darwin") {
		return "macOS"
	}
	if strings.Contains(ua, "linux") {
		return "Linux"
	}
	if strings.Contains(ua, "android") {
		return "Android"
	}
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod") {
		return "iOS"
	}

	return "Unknown"
}

func extractBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "edg/") {
		return "Microsoft Edge"
	}
	if strings.Contains(ua, "chrome/") {
		return "Chrome"
	}
	if strings.Contains(ua, "firefox/") {
		return "Firefox"
	}
	if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome") {
		return "Safari"
	}
	if strings.Contains(ua, "opera/") || strings.Contains(ua, "opr/") {
		return "Opera"
	}

	return "Unknown"
}
package utils

import (
	"net"
	"net/http"
	"strings"
)

// DeviceInfo represents extracted device information
type DeviceInfo struct {
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
	Platform  string `json:"platform"`
	Browser   string `json:"browser"`
}

// ExtractDeviceInfo extracts device information from an HTTP request
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

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (from proxy/load balancer)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if there are multiple
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// extractPlatform extracts the platform/OS from User-Agent
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

// extractBrowser extracts the browser from User-Agent
func extractBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	// Order matters - check more specific first
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
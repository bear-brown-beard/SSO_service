package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
)

func BrowserFingerprint(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	acceptLanguage := r.Header.Get("Accept-Language")
	acceptEncoding := r.Header.Get("Accept-Encoding")
	connection := r.Header.Get("Connection")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	fingerprintData := fmt.Sprintf("%s|%s|%s|%s|%s", userAgent, acceptLanguage, acceptEncoding, connection, ip)
	hash := sha256.Sum256([]byte(fingerprintData))
	fingerprint := hex.EncodeToString(hash[:])

	return fingerprint
}

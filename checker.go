package botchecker

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// list of private subnets.
var privateMasks = func(ips ...string) []net.IPNet {
	masks := make([]net.IPNet, 0, len(ips))
	for i := range ips {
		_, network, _ := net.ParseCIDR(ips[i])
		masks = append(masks, *network)
	}
	return masks
}("10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16")

// IsGoogleBot checks a request from GoogleBot or not.
func IsGoogleBot(r *http.Request) (bool, error) {

	if !strings.Contains(r.UserAgent(), "Googlebot") {
		return false, nil
	}

	ip := net.ParseIP(r.RemoteAddr)
	if ip == nil || !ip.IsGlobalUnicast() {
		return false, nil
	}

	for i := range privateMasks {
		if privateMasks[i].Contains(ip) {
			return false, nil
		}
	}

	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return false, err
	}

	host := fmt.Sprintf("crawl-%s.googlebot.com.", strings.Replace(ip.String(), ".", "-", -1))
	for i := range names {
		if host == names[i] {
			return true, nil
		}
	}

	return false, nil
}

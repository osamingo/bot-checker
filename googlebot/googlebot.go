package googlebot

import (
	"net"
	"net/http"
	"strings"

	"github.com/osamingo/bot-checker"
)

// BotTypeGooglebot type
const BotTypeGooglebot = botchecker.BotType("Googlebot")

// Checker struct.
type Checker struct{}

// NewGooglebotChecker returns googlebot.GooglebotCheker.
func NewGooglebotChecker() *Checker {
	return new(Checker)
}

// Check a request from GoogleBot or not.
func (c *Checker) Check(r *http.Request) (botchecker.BotType, error) {

	if !strings.Contains(r.UserAgent(), string(BotTypeGooglebot)) {
		return botchecker.BotTypeNoBot, nil
	}

	ip := net.ParseIP(r.RemoteAddr)
	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return botchecker.BotTypeNoBot, err
	}

	host := ""
	for i := range names {
		if strings.HasSuffix(names[i], ".googlebot.com.") || strings.HasSuffix(names[i], ".google.com.") {
			host = names[i]
			break
		}
	}

	if host == "" {
		return botchecker.BotTypeNoBot, nil
	}

	ret, err := net.LookupIP(host[:len(host)-1])
	if err != nil {
		return botchecker.BotTypeNoBot, err
	}

	for i := range ret {
		if ip.Equal(ret[i]) {
			return BotTypeGooglebot, nil
		}
	}

	return botchecker.BotTypeNoBot, nil
}

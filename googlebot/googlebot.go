package googlebot

import (
	"fmt"
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

	ip := r.RemoteAddr
	names, err := net.LookupAddr(ip)
	if err != nil {
		return botchecker.BotTypeNoBot, err
	}

	host := fmt.Sprintf("crawl-%s.googlebot.com.", strings.Replace(ip, ".", "-", 4))
	for i := range names {
		if host == names[i] {
			return BotTypeGooglebot, nil
		}
	}

	return botchecker.BotTypeNoBot, nil
}

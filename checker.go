package botchecker

import (
	"net"
	"net/http"
)

type (
	// A BotType is a type of bot.
	BotType string
	// BotChecker interface.
	BotChecker interface {
		Check(*http.Request) (BotType, error)
	}
)

// BotTypeNoBot is no bot request.
const BotTypeNoBot = BotType("")

var (
	// list of private subnets.
	privateMasks = toMasks("10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16")
	// defalut checkers.
	defaultCheckers []BotChecker
)

// converts a list of subnets' string to a list of net.IPNet.
func toMasks(ips ...string) []net.IPNet {
	masks := make([]net.IPNet, 0, len(ips))
	for i := range ips {
		_, network, _ := net.ParseCIDR(ips[i])
		masks = append(masks, *network)
	}
	return masks
}

// AddDefaultCheckers adds checker to default checkers.
func AddDefaultCheckers(checkers ...BotChecker) {
	for i := range checkers {
		defaultCheckers = append(defaultCheckers, checkers[i])
	}
}

// Do bot checks.
func Do(r *http.Request, checkers ...BotChecker) (BotType, error) {

	ua, addr := r.UserAgent(), r.RemoteAddr
	if ua == "" || addr == "" {
		return BotTypeNoBot, nil
	}

	ip := net.ParseIP(addr)
	if ip == nil || !ip.IsGlobalUnicast() {
		return BotTypeNoBot, nil
	}

	for i := range privateMasks {
		if privateMasks[i].Contains(ip) {
			return BotTypeNoBot, nil
		}
	}

	checkers = append(checkers, defaultCheckers...)

	for i := range checkers {
		bot, err := checkers[i].Check(r)
		if bot != BotTypeNoBot {
			return bot, nil
		}
		if err != nil {
			return BotTypeNoBot, err
		}
	}

	return BotTypeNoBot, nil
}

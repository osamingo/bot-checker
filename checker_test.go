package botchecker

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mock struct{}

func (m *mock) Check(r *http.Request) (BotType, error) {
	if r.UserAgent() == "fakerror" {
		return BotTypeNoBot, errors.New("fakerror")
	}

	if r.UserAgent() == "imbot" {
		return BotType("dummy"), nil
	}

	return BotTypeNoBot, nil
}

func TestIsGoogleBot(t *testing.T) {

	m := new(mock)

	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	bot, err := Do(r, m)
	require.NoError(t, err)
	assert.Equal(t, BotTypeNoBot, bot)

	r.Header.Set("User-Agent", "bot")
	r.RemoteAddr = "127.0.0.1"
	bot, err = Do(r, m)
	require.NoError(t, err)
	assert.Equal(t, BotTypeNoBot, bot)

	r.Header.Set("User-Agent", "Googlebot")
	bot, err = Do(r, m)
	require.NoError(t, err)
	assert.Equal(t, BotTypeNoBot, bot)

	ips := []string{
		// broadcast address
		"255.255.255.255",
		// private ips
		"10.0.0.1",
		"172.16.0.1",
		"192.168.0.0",
		// dev.abema.tv
		"130.211.15.86",
	}
	for _, ip := range ips {
		r.RemoteAddr = ip
		bot, err = Do(r, m)
		require.NoError(t, err)
		assert.Equal(t, BotTypeNoBot, bot)
	}

	r.RemoteAddr = "66.249.75.228"
	bot, err = Do(r, m)
	require.NoError(t, err)
	assert.Equal(t, BotTypeNoBot, bot)

	r.Header.Set("User-Agent", "fakerror")
	bot, err = Do(r, m)
	require.Error(t, err)
	assert.Equal(t, BotTypeNoBot, bot)

	r.Header.Set("User-Agent", "imbot")
	bot, err = Do(r, m)
	require.NoError(t, err)
	assert.Equal(t, BotType("dummy"), bot)
}

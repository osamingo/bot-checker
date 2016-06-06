package googlebot

import (
	"net/http"
	"testing"

	"github.com/osamingo/bot-checker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGooglebotCheker(t *testing.T) {

	gc := NewGooglebotChecker()
	require.NotNil(t, gc)

	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	bot, err := gc.Check(r)
	require.NoError(t, err)
	assert.Equal(t, botchecker.BotTypeNoBot, bot)

	r.Header.Set("User-Agent", "bot")
	bot, err = gc.Check(r)
	require.NoError(t, err)
	assert.Equal(t, botchecker.BotTypeNoBot, bot)

	r.Header.Set("User-Agent", "Googlebot")
	r.RemoteAddr = "127.0.0.1:23456"
	bot, err = gc.Check(r)
	require.NoError(t, err)
	assert.Equal(t, botchecker.BotTypeNoBot, bot)

	r.RemoteAddr = "255.255.255.200"
	_, err = gc.Check(r)
	require.Error(t, err)

	r.RemoteAddr = "66.249.75.228"
	bot, err = gc.Check(r)
	require.NoError(t, err)
	assert.Equal(t, BotTypeGooglebot, bot)
}

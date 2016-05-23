package botchecker

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGoogleBot(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	require.NoError(t, err)

	ret, err := IsGoogleBot(r)
	require.NoError(t, err)
	assert.False(t, ret)

	r.Header.Set("User-Agent", "Googlebot")
	ret, err = IsGoogleBot(r)
	require.NoError(t, err)
	assert.False(t, ret)

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
		ret, err = IsGoogleBot(r)
		require.NoError(t, err)
		assert.False(t, ret)
	}

	r.RemoteAddr = "255.255.255.200"
	ret, err = IsGoogleBot(r)
	require.Error(t, err)
	assert.False(t, ret)

	r.RemoteAddr = "66.249.75.228"
	ret, err = IsGoogleBot(r)
	require.NoError(t, err)
	assert.True(t, ret)
}

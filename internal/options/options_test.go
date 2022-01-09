package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetOreillyTrialOptions function tests if GetOreillyTrialOptions function running properly
func TestGetVarnishCacheInvalidatorOptions(t *testing.T) {
	t.Log("fetching default options.VarnishCacheInvalidatorOptions")
	opts := GetVarnishCacheInvalidatorOptions()
	assert.NotNil(t, opts)
	t.Logf("fetched default options.VarnishCacheInvalidatorOptions, %v\n", opts)
}

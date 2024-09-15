package httpclient_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudokyo/httpclient"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_E2E(t *testing.T) {
	t.Run("Test load default config", func(t *testing.T) {
		config := httpclient.LoadConfig("config")
		config.Print()

		assert.Equal(t, false, config.Debug)
		assert.Equal(t, "", config.Endpoint)
		assert.Equal(t, 3, config.RetryCount)
		assert.Equal(t, 30*time.Second, config.RetryTimeout)

		fmt.Println("--> OUT\t", config)
	})

	t.Run("Test load set config", func(t *testing.T) {
		debug := true
		retry := 1
		userAgent := "E2E Testing 1.x"
		endpoint := "https://www.google.com"
		timeout := 3 * time.Second

		config := httpclient.LoadConfig("config", debug, retry, timeout, endpoint, userAgent)
		config.Print()

		assert.Equal(t, debug, config.Debug)
		assert.Equal(t, endpoint, config.Endpoint)
		assert.Equal(t, retry, config.RetryCount)
		assert.Equal(t, timeout, config.RetryTimeout)
		assert.Equal(t, userAgent, config.UserAgent)

		fmt.Println("--> OUT\t", config)
	})

	t.Run("Test load config from ENV", func(t *testing.T) {
		// Set the env
		os.Setenv("CONFIG_RESTY_DEBUG", "true")
		os.Setenv("CONFIG_RESTY_RETRY_COUNT", "2")
		os.Setenv("CONFIG_RESTY_TIMEOUT", "3s")
		os.Setenv("CONFIG_RESTY_RETRY_WAIT_TIME", "1s")
		os.Setenv("CONFIG_RESTY_RETRY_MAX_WAIT_TIME", "2s")
		os.Setenv("CONFIG_RESTY_ENDPOINT", "https://www.google.com")

		config := httpclient.LoadConfig("config")
		config.Print()

		assert.Equal(t, true, config.Debug)
		assert.Equal(t, "https://www.google.com", config.Endpoint)
		assert.Equal(t, 2, config.RetryCount)
		assert.Equal(t, 3*time.Second, config.RetryTimeout)
		assert.Equal(t, 1*time.Second, config.RetryWaitTime)
		assert.Equal(t, 2*time.Second, config.RetryMaxWaitTime)

		fmt.Println("--> OUT\t", config)
	})
}

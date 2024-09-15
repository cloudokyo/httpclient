package httpclient_test

import (
	"fmt"
	"testing"

	"github.com/cloudokyo/httpclient"
	"github.com/stretchr/testify/assert"
)

func TestHttpClient_E2E(t *testing.T) {
	// Load config from ENV
	config := httpclient.LoadConfig("config", "HTTPCLIENT 1.x")
	config.Endpoint = "https://www.google.com"
	config.Print()

	// Init the resty client
	client := httpclient.NewResty(config)

	// Execute a request
	res, err := client.R().Get("https://www.google.com")

	assert.NoError(t, err)
	assert.NotNil(t, res)

	fmt.Println("--> OUT:", string(res.Body()))
}

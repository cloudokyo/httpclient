package httpclient

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudokyo/env"
	"github.com/go-resty/resty/v2"
)

// The config to create restry http client
type Config struct {
	Name             string
	Debug            bool
	Endpoint         string
	UserAgent        string
	ContentType      string
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	RetryTimeout     time.Duration
}

// The helper function converts the config to string
func (config Config) String() string {
	return fmt.Sprintf("%#v", config)
}

// The helper function print the data output
//
// Example
//
//	config := httpclient.LoadConfig()
//	config.Print()
func (config Config) Print() {
	log.Printf("- %s.debug: %v\n", config.Name, config.Debug)
	log.Printf("- %s.endpoint: %s\n", config.Name, config.Endpoint)
	log.Printf("- %s.retryCount: %d\n", config.Name, config.RetryCount)
	log.Printf("- %s.retryWaitTime: %s\n", config.Name, config.RetryWaitTime)
	log.Printf("- %s.retryMaxWaitTime: %s\n", config.Name, config.RetryMaxWaitTime)
	log.Printf("- %s.retryTimeout: %s\n", config.Name, config.RetryTimeout)
}

// Load the resty config from ENV
//
// Env
//
//	{NAME}_RESTY_{KEY}
//	CONFIG_RESTY_DEBUG
//
// Example
//
//	config := httpclient.LoadConfig("config")
//	config := httpclient.LoadConfig("config", 1) --> set {retryCount}
//	config := httpclient.LoadConfig("config", 1, true) --> set {retryCount, debug}
//	config := httpclient.LoadConfig("config", "https://www.google.com") --> set {endpoint}
//	config.Print()
func LoadConfig(name string, args ...any) Config {
	// Check the args
	var debug bool
	var retry int
	var userAgent string
	var endpoint string
	var timeout time.Duration

	for _, arg := range args {
		switch value := arg.(type) {
		case bool:
			debug = value
		case int:
			retry = value
		case string:
			if strings.HasPrefix(value, "http") {
				endpoint = value
			} else {
				userAgent = value
			}
		case time.Duration:
			timeout = value
		}
	}

	if retry == 0 {
		retry = 3
	}

	if timeout == 0 {
		timeout = 30 * time.Second
	}

	if userAgent == "" {
		userAgent = "RESTY " + resty.Version
	}

	prefix := strings.ToUpper(name)

	return Config{
		Name:             strings.ToLower(name),
		Debug:            env.GetBool(prefix+"_RESTY_DEBUG", debug),
		Endpoint:         env.Get(prefix+"_RESTY_ENDPOINT", endpoint),
		UserAgent:        env.Get(prefix+"_RESTY_USER_AGENT", userAgent),
		ContentType:      env.Get(prefix+"_RESTY_CONTENT_TYPE", "application/json"),
		RetryCount:       env.GetInt(prefix+"_RESTY_RETRY_COUNT", retry),
		RetryTimeout:     env.GetDuration(prefix+"_RESTY_TIMEOUT", timeout),
		RetryWaitTime:    env.GetDuration(prefix+"_RESTY_RETRY_WAIT_TIME", 5*time.Second),
		RetryMaxWaitTime: env.GetDuration(prefix+"_RESTY_RETRY_MAX_WAIT_TIME", 20*time.Second),
	}
}

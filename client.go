package httpclient

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/cloudokyo/env"
	"github.com/cloudokyo/log"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/http2"
)

var (
	http2Enable         = env.GetBool("HTTP_CLIENT_HTTP2", false)
	maxIdleConns        = env.GetInt("HTTP_CLIENT_MAX_IDLE", 100)
	maxIdleConnsPerHost = env.GetInt("HTTP_CLIENT_MAX_IDLE_PER_HOST", 100)
)

func init() {
	if http2Enable {
		log.Println("- http.client.protocol: HTTP/2.0")
	} else {
		log.Println("- http.client.protocol: HTTP/1.1")
		log.Println("- http.client.maxIdleConns:", maxIdleConns)
		log.Println("- http.client.maxIdleConnsPerHost:", maxIdleConnsPerHost)
	}
}

// Create a new http client that support larger connection pool or using HTTP/2.0
//
// We can change the http client by ENV config, see above section for more informations.
func New() *http.Client {
	// Check for HTTP/2.0
	if http2Enable {
		// Init the http.Client that works with http/2.0
		return &http.Client{Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP: true,

			// Pretend we are dialing a TLS endpoint. (we ignore the passed tls.Config)
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}}
	} else {
		// Customize the Transport to have larger connection pool
		roundTripper := http.DefaultTransport
		transportPointer, ok := roundTripper.(*http.Transport)
		if !ok {
			log.Fatal("roundTripper is not an *http.Transport")
		}

		transport := *transportPointer
		transport.MaxIdleConns = maxIdleConns
		transport.MaxIdleConnsPerHost = maxIdleConnsPerHost

		return &http.Client{Transport: &transport}
	}
}

// Create a new Resty http client using the config instance.
func NewResty(config Config) (client *resty.Client) {
	// Init a new resty client
	client = resty.NewWithClient(New())

	// Set debug mode
	//
	// Default: false
	client.SetDebug(config.Debug)

	// Set retry count to non zero to enable retries
	//
	// Default: 3
	client.SetRetryCount(config.RetryCount)

	// You can override initial retry wait time.
	//
	// Default: 5s
	client.SetRetryWaitTime(config.RetryWaitTime)

	// MaxWaitTime can be overridden as well.
	//
	// Default: 20s
	client.SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	// Set client timeout as per your need
	//
	// Default: 30s
	client.SetTimeout(config.RetryTimeout)

	// Host URL for all request. So you can use relative URL in the request
	client.SetBaseURL(config.Endpoint)

	// Headers for all request
	client.SetHeaders(map[string]string{
		"Content-Type": config.ContentType,
		"User-Agent":   config.UserAgent,
	})

	return client
}

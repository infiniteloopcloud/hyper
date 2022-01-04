package hyper

import (
	"net/http"
	"time"
)

var (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
	timeout             = 10 * time.Second
)

func Client() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost

	return &http.Client{
		Timeout:   timeout,
		Transport: t,
	}
}

func SetClientMaxIdleConns(v int)        { maxIdleConns = v }
func SetClientMaxConnsPerHost(v int)     { maxConnsPerHost = v }
func SetClientMaxIdleConnsPerHost(v int) { maxIdleConnsPerHost = v }
func SetClientTimeout(v time.Duration)   { timeout = v }

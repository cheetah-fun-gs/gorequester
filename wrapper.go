package gorequester

import (
	"time"
)

func defaultClient() *Client {
	return &Client{
		Timeout: time.Second * 5,
	}
}

// New New
func New(method, toURL string) *Requester {
	return defaultClient().New(method, toURL)
}

// Get Get
func Get(toURL string) *Requester {
	return defaultClient().Get(toURL)
}

// Post Post
func Post(toURL string) *Requester {
	return defaultClient().Post(toURL)
}

// PostData PostData
func PostData(toURL, contentType string, v interface{}) *Requester {
	return defaultClient().PostData(toURL, contentType, v)
}

// PostJSON PostJSON
func PostJSON(toURL string, v interface{}) *Requester {
	return defaultClient().PostJSON(toURL, v)
}

package langpractice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/ebitengine/oto/v3"
	"net/http"
	"time"
)

const (
	BaseHREF = `https://langpractice.com/api/newnumber/spanish-mexico/1`
)

type LangPractice struct {
	client      *http.Client
	url         string
	PlayTimeout time.Duration
	OtoContext  *oto.Context

	// AutoNext if true, play next number automatically or wait on any keypress
	AutoNext bool
}

// RequestNumber makes a request for a new random number from langpractice.com and returns the unmarshalled response
func (c *LangPractice) RequestNumber() (*LPResponse, error) {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return c.parseResponse(resp.Body)
}

// NewClient creates a new instance of the client
func NewClient(min, max int) *LangPractice {
	c := LangPractice{
		PlayTimeout: time.Second * 10, // time to play a spanish number audio
	}

	// Create custom Transport with TLS configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // Set to true if self-signed certificates are used (not recommended)
		},
	}

	// Create an HTTP client with timeout and custom transport
	c.client = &http.Client{
		Timeout:   10 * time.Second, // Set a timeout for requests
		Transport: tr,
	}

	// Make URL for a GET request
	c.url = fmt.Sprintf("%s/%d/%d", BaseHREF, min, max)

	// new audio context
	ctx, err := NewPlayer()
	if err != nil {
		return nil
	}
	c.OtoContext = ctx

	return &c
}

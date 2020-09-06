package rest

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/thejasn/tester/core/client"
)

type Config struct {
	ctx     context.Context
	client  *http.Client
	request *http.Request
	headers map[string]string
	method  string
	body    string
	baseURL string
	url     string
}

func NewRestConfig(ctx context.Context, baseURL string) *Config {
	return &Config{
		ctx:     ctx,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func WithHeaders(headers map[string]string) client.RunnerOpts {
	return func(p client.Runner) {
		p.(*Config).headers = headers
	}
}

func WithBody(body string) client.RunnerOpts {
	return func(p client.Runner) {
		p.(*Config).body = body
	}
}

func WithMethod(method string) client.RunnerOpts {
	return func(p client.Runner) {
		p.(*Config).method = method
	}
}

func WithUriPath(url string) client.RunnerOpts {
	return func(p client.Runner) {
		p.(*Config).url = url
	}
}

func (c Config) GetIdentifier() string {
	return "1"
}

func (c *Config) Build() error {
	var err error
	var req *http.Request
	if c.method == "POST" {
		req, err = http.NewRequestWithContext(c.ctx, c.method, strings.Join([]string{c.baseURL, c.url}, ""), strings.NewReader(c.body))
	} else if c.method == "GET" {
		req, err = http.NewRequestWithContext(c.ctx, c.method, strings.Join([]string{c.baseURL, c.url}, ""), nil)
	} else {
		return errors.New("unsupported method")
	}
	if err != nil {
		return err
	}
	c.request = req
	for k, v := range c.headers {
		c.request.Header.Add(k, v)
	}
	return nil
}

func (c *Config) Invoke() (string, error) {
	resp, err := c.client.Do(c.request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Config) Clear() {
	c.body = ""
	c.headers = nil
	c.method = ""
	c.url = ""
	c.request = nil
}

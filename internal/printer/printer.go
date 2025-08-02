package printer

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Printer struct {
	client *http.Client
}

func New() *Printer {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Printer{
		client: client,
	}
}

func (p *Printer) Cut() error {
	_, err := p.client.Post("http://10.100.0.3:8000/cut", "application/x-www-form-urlencoded", nil)
	return err
}

func (p *Printer) Text(s string) error {
	data := url.Values{}
	data.Set("text", s)
	data.Set("cut", "true")

	_, err := p.client.Post("http://10.100.0.3:8000/text", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	return err
}

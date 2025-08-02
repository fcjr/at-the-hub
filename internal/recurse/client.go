package recurse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://www.recurse.com"

type Client struct {
	token  string
	client *http.Client
}

func NewClient(token string) (*Client, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Client{
		token:  token,
		client: client,
	}, nil
}

func (c *Client) buildRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	return req, nil
}

func (c *Client) Profile(id int) (*Profile, error) {
	req, err := c.buildRequest("GET", fmt.Sprintf("/api/v1/profiles/%d", id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile Profile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (c *Client) CurrentCheckins() ([]Checkin, error) {
	req, err := c.buildRequest("GET", fmt.Sprintf("/api/v1/hub_visits?date=%s", time.Now().Format("2006-01-02")), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var checkins []Checkin
	err = json.NewDecoder(resp.Body).Decode(&checkins)
	if err != nil {
		return nil, err
	}

	return checkins, nil
}

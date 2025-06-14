package loggerclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Log struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}
}

func (c *Client) Log(ctx context.Context, l Log) error {
	jsonData, _ := json.MarshalIndent(l, "", "\t")
	u := fmt.Sprintf("%s/log", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("logger returned %s", resp.Status)
	}
	return nil
}

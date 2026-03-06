package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Message      MessageBody `json:"message"`
	Subscription string      `json:"subscription"`
}

type MessageBody struct {
	Data       string            `json:"data"` // base64-encoded
	MessageID  string            `json:"messageId"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type Pusher struct {
	Endpoint string
	Client   *http.Client
}

func New(endpoint string) *Pusher {
	return &Pusher{
		Endpoint: endpoint,
		Client:   http.DefaultClient,
	}
}

func (p *Pusher) Push(msg *Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", p.Endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request to '%s' failed with status '%s'", p.Endpoint, resp.Status)
	}
	return nil
}

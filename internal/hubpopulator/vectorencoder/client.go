package vectorencoder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type EncodeRequest struct {
	Input Input `json:"input"`
}

type Input struct {
	Text string `json:"text"`
}

type EncodeResponse struct {
	//Input struct {
	//	Text string `json:"text"`
	//} `json:"input"`
	Output []float64 `json:"output"`
	//ID                  interface{} `json:"id"`
	CreatedAt   interface{} `json:"created_at"`
	StartedAt   time.Time   `json:"started_at"`
	CompletedAt time.Time   `json:"completed_at"`
	Logs        string      `json:"logs"`
	//Error               /interface{} `json:"error"`
	Status string `json:"status"`
	//WebhookEventsFilter []string    `json:"webhook_events_filter"`
	//OutputFilePrefix    interface{} `json:"output_file_prefix"`
	//Webhook             interface{} `json:"webhook"`
}

type Client struct {
	encoderURL string
	client     *http.Client
}

func NewClient(encoderURL string) *Client {
	return &Client{
		encoderURL: encoderURL,
		client:     &http.Client{},
	}
}

func (c *Client) Encode(text string) (*EncodeResponse, error) {
	encodeStruct := EncodeRequest{Input: Input{Text: text}}
	body, err := json.Marshal(&encodeStruct)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.encoderURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	response := &EncodeResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

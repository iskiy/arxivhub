package pinecone

import (
	"arxivhub/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	insertEndpoint = "/vectors/upsert"
	queryEndpoint  = "/query"
	fetchEndpoint  = "/vectors/fetch"
)

type InsertRequest struct {
	Vector []Vector `json:"vectors"`
}

type Vector struct {
	Values []float64 `json:"values"`
	ID     string    `json:"id"`
}

type Client struct {
	apiKey string
	host   string
	client *http.Client
}

func NewClient(apiKey string, host string) *Client {
	client := http.Client{}
	return &Client{
		apiKey: apiKey,
		host:   host,
		client: &client,
	}
}

func (c *Client) InsertVectors(data []byte) error {
	path := "https://" + c.host + insertEndpoint
	req, err := http.NewRequest("POST", path, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Host", c.host)
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	fmt.Println(resp.StatusCode)
	return err
}

type QueryRequest struct {
	IncludeValues   bool      `json:"includeValues"`
	IncludeMetadata bool      `json:"includeMetadata"`
	Vector          []float64 `json:"vector"`
	TopK            int       `json:"topK"`
	Filter          Filter    `json:"filter"`
}

type Filter struct {
	Year     int    `json:"year,omitempty"`
	Category string `json:"categories,omitempty"`
}

type QueryResponse struct {
	Results []any `json:"results"`
	Matches []struct {
		ID        string           `json:"id"`
		Score     float64          `json:"score"`
		PaperData models.PaperData `json:"metadata"`
	} `json:"matches"`
}

func (c *Client) Query(request QueryRequest) (*QueryResponse, error) {
	path := "https://" + c.host + queryEndpoint

	data, err := json.Marshal(request)
	if err != nil {
		return &QueryResponse{}, err
	}
	req, err := http.NewRequest("POST", path, bytes.NewReader(data))
	if err != nil {
		return &QueryResponse{}, err
	}

	req.Header.Set("Host", c.host)
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	queryResponse := QueryResponse{}
	err = json.NewDecoder(resp.Body).Decode(&queryResponse)
	if err != nil {
		return &QueryResponse{}, err
	}

	return &queryResponse, err
}

type PaperRecord struct {
	ID        string           `json:"id"`
	PaperData models.PaperData `json:"metadata"`
}

type PaperRecords map[string]PaperRecord

type FetchResponse struct {
	Vectors map[string]PaperRecord `json:"vectors"`
}

func (c *Client) FetchPapers(ids []string) (*FetchResponse, error) {
	query := "?"

	for _, id := range ids {
		query += "ids=" + id + "&"
	}
	path := "https://" + c.host + fetchEndpoint + query
	fmt.Println(path)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Host", c.host)
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(req)

	fetchResponse := FetchResponse{}
	err = json.NewDecoder(response.Body).Decode(&fetchResponse)
	if err != nil {
		return nil, fmt.Errorf("parse pinecode response error: %v", err)
	}
	return &fetchResponse, err
}

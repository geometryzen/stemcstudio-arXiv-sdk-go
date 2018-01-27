package arXiv

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// SearchRequest is ...
type SearchRequest struct {
	Query string `json:"query"`
}

// Submission is ...
type Submission struct {
	HRef     string   `json:"href"`
	Owner    string   `json:"owner"`
	GistID   string   `json:"gistId"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Keywords []string `json:"keywords"`
}

// SearchResponse is ...
type SearchResponse struct {
	Found int64        `json:"found"`
	Start int64        `json:"start"`
	Refs  []Submission `json:"refs,omitempty"`
}

// Service is ...
type Service interface {
	Search(query string, size int) (*SearchResponse, error)
	Submit(packet *Submission) (interface{}, error)
}

// Client probably does not need to be exposed.
type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

// NewClient is a factory for creating a search service.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{httpClient: httpClient}
}

// Search is ...
func (c *Client) Search(query string, size int) (*SearchResponse, error) {
	rel := &url.URL{Path: "/search"}
	u := c.BaseURL.ResolveReference(rel)
	data := SearchRequest{Query: query}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bundle *SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&bundle)
	return bundle, err
}

// Submit is ...
func (c *Client) Submit(packet *Submission) (interface{}, error) {
	rel := &url.URL{Path: "/submit"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// var bundle *searchResponse
	// err = json.NewDecoder(resp.Body).Decode(&bundle)
	return nil, err
}

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Request struct represents an HTTP request to be sent to the Azure OpenAI API
type Request struct {
	body    any
	method  string
	url     string
	headers map[string]interface{}
}

// It contains the request body, method, URL, and headers.
// The body can be of any type, and the headers are a map of string keys to interface{} values.
// The method is the HTTP method (e.g., GET, POST) to be used for the request.
// The URL is the endpoint of the Azure OpenAI API to which the request will be sent.
// The headers are optional and can be used to set any additional headers required by the API.
func NewRequest(method string, url string, body any, headers map[string]interface{}) *Request {
	return &Request{
		body:    body,
		method:  method,
		url:     url,
		headers: headers,
	}
}

// The Request struct is used to create and send requests to the Azure OpenAI API.
// It has a method Send() that sends the request and returns the response.
func (r *Request) Send() (*http.Response, error) {
	data, err := json.Marshal(r.body)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return nil, err
	}
	client := &http.Client{}
	method := strings.ToUpper(r.method)
	req, err := http.NewRequest(method, r.url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	for key, value := range r.headers {
		req.Header.Set(key, value.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Azure AI Client of the Azure OpenAI API
// This client is used to send requests to the Azure OpenAI API
type AzureAIClient struct {
	// make the data private and add getters
	endpoint string
	apiKey   string
}

// and receive responses. It is initialized with the endpoint and API key.
func NewAzureAIClient(endpoint string, apiKey string) *AzureAIClient {
	return &AzureAIClient{
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

// Getters for the AzureAIClient
func (c *AzureAIClient) GetEndpoint() string {
	return c.endpoint
}
func (c *AzureAIClient) GetApiKey() string {
	return c.apiKey
}

func (c *AzureAIClient) CreateCompletions(prompt string) (any, error) {
	body := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "You are an experience software engineer, with more that 20 years experience with a specialization in backend development with vast knowledge in Python, Go, CPP, Java, Rust, Typescript and SQL, you mostly be helping in transforming raw schema to an ORM model, for Python use Django orm/ Schalchemy if specified, for Go use GORM, for rust use diesel, for Typescript use Typeorm, for java use springboot jdbc.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 1,
		"top_p":       1,
		"model":       "openai/gpt-4.1",
	}
	// Create a new request
	req := NewRequest("POST", c.endpoint, body, map[string]interface{}{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	})
	response, err := req.Send()
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	// Read the response body
	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}
	// Return the response body
	return data, nil

}

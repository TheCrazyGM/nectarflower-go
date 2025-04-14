package hive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client represents a Hive API client
type Client struct {
	Nodes        []string
	FailingNodes map[string]string
	httpClient   *http.Client
}

// NewClient creates a new Hive client with default node
func NewClient() *Client {
	return &Client{
		Nodes:        []string{"https://api.hive.blog"},
		FailingNodes: make(map[string]string),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SetNodes updates the client's node list
func (c *Client) SetNodes(nodes []string, failingNodes map[string]string) {
	validNodes := make([]string, 0, len(nodes))

	// Filter out invalid URLs and nodes that are in the failing_nodes list
	for _, node := range nodes {
		// Skip if node is in failing_nodes list
		if _, isFailing := failingNodes[node]; isFailing {
			continue
		}

		// Validate URL format
		_, err := url.Parse(node)
		if err == nil {
			validNodes = append(validNodes, node)
		}
	}

	if len(validNodes) > 0 {
		c.Nodes = validNodes
	}
	c.FailingNodes = failingNodes
}

// RPCRequest represents a JSON-RPC request
type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

// RPCResponse represents a JSON-RPC response
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Call makes a JSON-RPC call to the Hive API
func (c *Client) Call(method string, params any, result any) error {
	// Try each node until one succeeds
	var lastErr error
	for _, node := range c.Nodes {
		err := c.callNode(node, method, params, result)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return fmt.Errorf("all nodes failed: %w", lastErr)
	}

	return errors.New("no nodes available")
}

// callNode makes a JSON-RPC call to a specific node
func (c *Client) callNode(node string, method string, params any, result any) error {
	// Create request
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	// Marshal request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", node, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	var rpcResponse RPCResponse
	err = json.NewDecoder(resp.Body).Decode(&rpcResponse)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	// Check for RPC error
	if rpcResponse.Error != nil {
		return fmt.Errorf("RPC error: %s (code: %d)", rpcResponse.Error.Message, rpcResponse.Error.Code)
	}

	// Unmarshal result
	err = json.Unmarshal(rpcResponse.Result, result)
	if err != nil {
		return fmt.Errorf("error unmarshaling result: %w", err)
	}

	return nil
}

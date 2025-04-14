package hive

import (
	"encoding/json"
	"fmt"
)

// AccountParams represents the parameters for the find_accounts method
type AccountParams struct {
	Accounts []string `json:"accounts"`
}

// Account represents a Hive account
type Account struct {
	Name         string `json:"name"`
	JSONMetadata string `json:"json_metadata"`
}

// AccountsResponse represents the response from find_accounts
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// NodeData represents the node information extracted from account metadata
type NodeData struct {
	Nodes        []string          `json:"nodes"`
	FailingNodes map[string]string `json:"failing_nodes"`
}

// GetNodesFromAccount fetches account JSON metadata and extracts node information
func (c *Client) GetNodesFromAccount(accountName string) (*NodeData, error) {
	// Create parameters for find_accounts method
	params := AccountParams{
		Accounts: []string{accountName},
	}

	// Make API call to get account information
	var response AccountsResponse
	err := c.Call("database_api.find_accounts", params, &response)
	if err != nil {
		return nil, fmt.Errorf("error fetching account: %w", err)
	}

	// Check if account exists
	if len(response.Accounts) == 0 {
		return nil, fmt.Errorf("account '%s' not found", accountName)
	}

	// Get account JSON metadata
	account := response.Accounts[0]
	jsonMetadata := account.JSONMetadata

	// Parse JSON metadata
	var metadataObj map[string]json.RawMessage
	err = json.Unmarshal([]byte(jsonMetadata), &metadataObj)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON metadata: %w", err)
	}

	// Extract nodes from metadata
	nodeData := &NodeData{
		Nodes:        make([]string, 0),
		FailingNodes: make(map[string]string),
	}

	// Extract nodes array
	if nodesJSON, ok := metadataObj["nodes"]; ok {
		err = json.Unmarshal(nodesJSON, &nodeData.Nodes)
		if err != nil {
			return nil, fmt.Errorf("error parsing nodes: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no nodes found in account metadata")
	}

	// Extract failing_nodes map if it exists
	if failingNodesJSON, ok := metadataObj["failing_nodes"]; ok {
		err = json.Unmarshal(failingNodesJSON, &nodeData.FailingNodes)
		if err != nil {
			// Just log the error but continue, failing_nodes is optional
			fmt.Printf("Warning: error parsing failing_nodes: %v\n", err)
		}
	}

	return nodeData, nil
}

// UpdateNodesFromAccount fetches nodes from an account and updates the client
func (c *Client) UpdateNodesFromAccount(accountName string) error {
	// Get nodes from account
	nodeData, err := c.GetNodesFromAccount(accountName)
	if err != nil {
		return err
	}

	// Update client with new nodes
	c.SetNodes(nodeData.Nodes, nodeData.FailingNodes)
	return nil
}

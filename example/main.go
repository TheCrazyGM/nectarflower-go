package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/thecrazygm/nectarflower-go/hive"
)

func main() {
	// Create a new Hive client with default node
	client := hive.NewClient()
	fmt.Println("Default client initialized with:", client.Nodes)

	// Account to fetch nodes from
	accountName := "nectarflower"

	// Get nodes from account
	fmt.Printf("\nFetching nodes from account %s...\n", accountName)
	nodeData, err := client.GetNodesFromAccount(accountName)
	if err != nil {
		log.Fatalf("Error fetching nodes: %v", err)
	}

	// Print node information
	fmt.Printf("Found %d nodes in account metadata\n", len(nodeData.Nodes))
	fmt.Println("Nodes:", nodeData.Nodes)

	if len(nodeData.FailingNodes) > 0 {
		fmt.Printf("Found %d failing nodes in account metadata\n", len(nodeData.FailingNodes))
		failingNodesJSON, _ := json.MarshalIndent(nodeData.FailingNodes, "", "  ")
		fmt.Printf("Failing nodes: %s\n", string(failingNodesJSON))
	}

	// Update client with new nodes
	fmt.Println("\nUpdating client with new nodes...")
	client.SetNodes(nodeData.Nodes, nodeData.FailingNodes)
	fmt.Println("Updated client initialized with:", client.Nodes)

	// Test the updated client with a simple query
	fmt.Println("\nTesting updated client with a query...")
	var props map[string]any
	err = client.Call("database_api.get_dynamic_global_properties", nil, &props)
	if err != nil {
		log.Fatalf("Error fetching global properties: %v", err)
	}

	// Convert the block number to int64 to avoid scientific notation
	blockNum, ok := props["head_block_number"].(float64)
	if ok {
		fmt.Printf("Query successful! Current block number: %d\n", int64(blockNum))
	} else {
		fmt.Printf("Query successful! Current block number: %v\n", props["head_block_number"])
	}

	// Demonstrate the all-in-one function
	fmt.Println("\nDemonstrating the all-in-one UpdateNodesFromAccount function...")
	newClient := hive.NewClient()
	err = newClient.UpdateNodesFromAccount(accountName)
	if err != nil {
		log.Fatalf("Error updating nodes: %v", err)
	}

	fmt.Println("One-step update complete. Client initialized with:", newClient.Nodes)

	// Example: Fetch a recent block
	fmt.Println("\nFetching a recent block...")
	
	// First get the current block number
	var blockProps map[string]any
	err = client.Call("database_api.get_dynamic_global_properties", nil, &blockProps)
	if err != nil {
		log.Fatalf("Error fetching global properties: %v", err)
	}
	
	currentBlockNum, ok := blockProps["head_block_number"].(float64)
	if !ok {
		log.Fatalf("Error getting block number")
	}
	
	// Fetch a block that's a few blocks behind the head to ensure it's available
	targetBlockNum := int64(currentBlockNum) - 10
	fmt.Printf("Fetching block #%d...\n", targetBlockNum)
	
	// Create parameters for get_block method
	blockParams := map[string]interface{}{
		"block_num": targetBlockNum,
	}
	
	// Fetch the block
	var block map[string]any
	err = client.Call("block_api.get_block", blockParams, &block)
	if err != nil {
		log.Fatalf("Error fetching block: %v", err)
	}
	
	// Extract block data
	blockData, ok := block["block"].(map[string]any)
	if !ok {
		log.Fatalf("Error extracting block data")
	}
	
	// Print block details
	fmt.Println("Block details:")
	fmt.Printf("  Block ID: %v\n", blockData["block_id"])
	fmt.Printf("  Previous: %v\n", blockData["previous"])
	fmt.Printf("  Timestamp: %v\n", blockData["timestamp"])
	
	// Print transaction count
	transactions, ok := blockData["transactions"].([]any)
	if ok {
		fmt.Printf("  Transaction count: %d\n", len(transactions))
		
		// If there are transactions, print details of the first one
		if len(transactions) > 0 {
			tx := transactions[0].(map[string]any)
			txId, _ := tx["transaction_id"].(string)
			fmt.Println("\nFirst transaction details:")
			fmt.Printf("  Transaction ID: %s\n", txId)
			
			// Pretty print the first transaction
			txJSON, _ := json.MarshalIndent(tx, "  ", "  ")
			fmt.Printf("  Transaction data:\n%s\n", string(txJSON))
		}
	}
	
	fmt.Println("\nDemonstration completed successfully!")
}

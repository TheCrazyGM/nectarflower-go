# nectarflower-go

A simple Golang library for interacting with Hive blockchain nodes, similar to nectarflower-js. This library allows you to fetch node information from a Hive account's JSON metadata and make API calls to the Hive blockchain.

## Features

- Fetch node information from a Hive account's JSON metadata
- Filter out failing nodes
- Make JSON-RPC calls to Hive API endpoints
- Automatically retry failed calls on different nodes

## Installation

```bash
go get github.com/thecrazygm/nectarflower-go
```

## Usage

### Simplest Use Case: Just Get Passing Nodes

If you just need a list of passing Hive nodes for your own client implementation:

```go
package main

import (
    "fmt"
    "log"

    "github.com/thecrazygm/nectarflower-go/hive"
)

func main() {
    // Create a client just to fetch nodes
    client := hive.NewClient()

    // Get nodes from account
    nodeData, err := client.GetNodesFromAccount("nectarflower")
    if err != nil {
        log.Fatalf("Error fetching nodes: %v", err)
    }

    // Print the list of passing nodes
    fmt.Println("Passing nodes:")
    for _, node := range nodeData.Nodes {
        fmt.Println(node)
    }
}
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/thecrazygm/nectarflower-go/hive"
)

func main() {
    // Create a new client with default node
    client := hive.NewClient()

    // Update nodes from account
    err := client.UpdateNodesFromAccount("nectarflower")
    if err != nil {
        log.Fatalf("Error updating nodes: %v", err)
    }

    // Make API calls using the updated client
    var props map[string]any
    err = client.Call("database_api.get_dynamic_global_properties", nil, &props)
    if err != nil {
        log.Fatalf("Error fetching global properties: %v", err)
    }

    // Convert the block number to int64 to avoid scientific notation
    blockNum, ok := props["head_block_number"].(float64)
    if ok {
        fmt.Printf("Current block number: %d\n", int64(blockNum))
    } else {
        fmt.Printf("Current block number: %v\n", props["head_block_number"])
    }
}
```

### Advanced Usage

```go
// Get nodes from account without updating the client
nodeData, err := client.GetNodesFromAccount("nectarflower")
if err != nil {
    log.Fatalf("Error fetching nodes: %v", err)
}

// Manually set nodes
client.SetNodes(nodeData.Nodes, nodeData.FailingNodes)

// Make a custom API call
params := map[string]any{
    "account": "nectarflower",
}

var result map[string]any
err = client.Call("condenser_api.get_accounts", []any{[]string{"nectarflower"}}, &result)
```

### Fetching Block Data

```go
// Create a client
client := hive.NewClient()

// Get the current block number
var props map[string]any
err := client.Call("database_api.get_dynamic_global_properties", nil, &props)
if err != nil {
    log.Fatalf("Error fetching global properties: %v", err)
}

// Convert to int64
currentBlockNum, ok := props["head_block_number"].(float64)
if !ok {
    log.Fatalf("Error getting block number")
}

// Fetch a specific block
targetBlockNum := int64(currentBlockNum) - 10
blockParams := map[string]any{
    "block_num": targetBlockNum,
}

var block map[string]any
err = client.Call("block_api.get_block", blockParams, &block)
if err != nil {
    log.Fatalf("Error fetching block: %v", err)
}

// Access block data
blockData, ok := block["block"].(map[string]any)
if ok {
    fmt.Printf("Block ID: %v\n", blockData["block_id"])
}
```

## Example

See the `example/main.go` file for a complete example of how to use the library.

## License

MIT

package graph

import (
	"fmt"
	"sync"

	"github.com/cheekybits/genny/generic"
)

type Item generic.Type

type Node struct {
	value Item
}

func (node *Node) String() string {
	return fmt.Sprintf("%v", node.value)
}

type ItemGraph struct {
	nodes       []*Node          // Consider doing a map: city_name -> *Node
	connections map[Node][]*Node // map: Node -> [*Node] (NWSE)
	lock        sync.RWMutex
}

func (graph *ItemGraph) AddNode(node *Node) {
	graph.lock.Lock()
	graph.nodes = append(graph.nodes, node)
	graph.lock.Unlock()
}

func (graph *ItemGraph) AddConnection(n1 *Node, n2 *Node) {
	graph.lock.Lock()
	if graph.connections == nil {
		graph.connections = make(map[Node][]*Node)
	}
	graph.connections[*n1] = append(graph.connections[*n1], n2)
	graph.connections[*n2] = append(graph.connections[*n2], n1)
	graph.lock.Unlock()
}

func (graph *ItemGraph) RemoveNode(value string) {
	graph.lock.Lock()
	// Remove the from the slice of Nodes
	for i := range graph.nodes {
		val := graph.nodes[i].String()
		if value == val {
			copy(graph.nodes[i:], graph.nodes[i+1:])
			graph.nodes[len(graph.nodes)-1] = nil
			graph.nodes = graph.nodes[:len(graph.nodes)-1]
			break
		}
	}

	// Check for any connections and remove
	graph.lock.Unlock()
}

func (graph *ItemGraph) RemoveConnections(value string) {
	graph.lock.Lock()
	// Look up our node in the map and receive the slice
	// That shows the connected nodes
	// Look up the connected node's links and remove the
	// current node from all connected nodes
	// Delete the key in the map and ensure the node gets
	// garbage collected
	graph.lock.Unlock()
}

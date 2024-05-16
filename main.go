package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Room struct represents a room in the ant farm
type Room struct {
	Name     string
	X        int
	Y        int
	AntCount int
}

// Tunnel struct represents a tunnel connecting two rooms
type Tunnel struct {
	From string
	To   string
	Used bool
}

// Graph struct represents an undirected graph
type Graph struct {
	Nodes map[string][]string
}

// NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string][]string),
	}
}

// AddEdge adds an undirected edge between two nodes
func (g *Graph) AddEdge(node1, node2 string) {
	g.Nodes[node1] = append(g.Nodes[node1], node2)
	g.Nodes[node2] = append(g.Nodes[node2], node1)
}

// Helper function to convert string to int
func convertToInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
func (g *Graph) PrintGraph() {
	for node, neighbors := range g.Nodes {
		for _, neighbor := range neighbors {
			fmt.Printf("%s-%s\n", node, neighbor)
		}
	}
}

// ReadFile function reads the input file and processes the rooms and tunnels
func ReadFile(filename string) (int, []Room, []Tunnel, Room, Room, error) {
	var antCount int
	var rooms []Room
	var tunnels []Tunnel
	var startRoom, endRoom Room

	file, err := os.Open(filename)
	if err != nil {
		return 0, nil, nil, Room{}, Room{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##start") {
			scanner.Scan()
			parts := strings.Fields(scanner.Text())
			startRoom = Room{
				Name: parts[0],
				X:    convertToInt(parts[1]),
				Y:    convertToInt(parts[2]),
			}
			rooms = append(rooms, startRoom) // Add start room to rooms slice
		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			parts := strings.Fields(scanner.Text())
			endRoom = Room{
				Name: parts[0],
				X:    convertToInt(parts[1]),
				Y:    convertToInt(parts[2]),
			}
			rooms = append(rooms, endRoom) // Add end room to rooms slice
		} else if antCount == 0 {
			antCount = convertToInt(line)
		} else if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			tunnels = append(tunnels, Tunnel{From: parts[0], To: parts[1]})
		} else {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				room := Room{
					Name: parts[0],
					X:    convertToInt(parts[1]),
					Y:    convertToInt(parts[2]),
				}
				rooms = append(rooms, room)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, nil, nil, Room{}, Room{}, err
	}

	return antCount, rooms, tunnels, startRoom, endRoom, nil
}
func (g *Graph) DFSWithStartAndEnd(startNode, endNode string, visited map[string]bool) {
	if visited == nil {
		visited = make(map[string]bool)
	}

	stack := []string{startNode}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if node == endNode {
			fmt.Println(node)
			break
		}

		if !visited[node] {
			fmt.Println(node)
			visited[node] = true

			for _, neighbor := range g.Nodes[node] {
				if !visited[neighbor] {
					stack = append(stack, neighbor)
				}
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . file.txt")
		os.Exit(1)
	}

	filename := os.Args[1]
	antCount, rooms, tunnels, startRoom, endRoom, err := ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Create a new graph
	graph := NewGraph()

	// Add edges to the graph
	for _, tunnel := range tunnels {
		graph.AddEdge(tunnel.From, tunnel.To)
	}

	// Print number of ants
	fmt.Printf("%d\n", antCount)

	// Print rooms
	fmt.Println("the_rooms")
	for _, room := range rooms {
		fmt.Printf("%s %d %d\n", room.Name, room.X, room.Y)
	}

	// Print links
	fmt.Println("the_links")
	for node, neighbors := range graph.Nodes {
		for _, neighbor := range neighbors {
			fmt.Printf("%s-%s\n", node, neighbor)
		}
	}
	fmt.Println("the_graph")
	graph.PrintGraph()

	fmt.Println("DFS Traversal from start to end:")
	graph.DFSWithStartAndEnd(startRoom.Name, endRoom.Name, nil)

}

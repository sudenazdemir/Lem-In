package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Room struct represents a room in the ant farm
type Room struct {
	Name string
	X    int
	Y    int
}

// Tunnel struct represents a tunnel connecting two rooms
type Tunnel struct {
	From string
	To   string
}

// Path struct represents a path taken by an ant
type Path struct {
	Steps []string
}

// Graph struct represents the graph structure
type Graph struct {
	Nodes    map[string]*Room
	Tunnels  map[string][]string
	Capacity map[string]map[string]int
	Start    string
	End      string
}

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
		} else if strings.HasPrefix(line, "##end") {
			scanner.Scan()
			parts := strings.Fields(scanner.Text())
			endRoom = Room{
				Name: parts[0],
				X:    convertToInt(parts[1]),
				Y:    convertToInt(parts[2]),
			}
		} else if antCount == 0 {
			antCount = convertToInt(line)
			if antCount <= 0 {
				return 0, nil, nil, Room{}, Room{}, fmt.Errorf("ERROR: No ants specified or invalid number of ants")
			}
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

// Helper function to convert string to int
func convertToInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// NewGraph initializes a new Graph
func NewGraph(antCount int, rooms []Room, tunnels []Tunnel, startRoom Room, endRoom Room) (*Graph, error) {
	graph := &Graph{
		Nodes:    make(map[string]*Room),
		Tunnels:  make(map[string][]string),
		Capacity: make(map[string]map[string]int),
		Start:    startRoom.Name,
		End:      endRoom.Name,
	}

	// Add start and end rooms
	graph.Nodes[startRoom.Name] = &startRoom
	graph.Nodes[endRoom.Name] = &endRoom

	// Add other rooms
	for i := range rooms {
		room := rooms[i]
		graph.Nodes[room.Name] = &room
	}

	for _, tunnel := range tunnels {
		if tunnel.From == tunnel.To {
			return nil, fmt.Errorf("ERROR: Room %s links to itself", tunnel.From)
		}

		if _, ok := graph.Tunnels[tunnel.From]; !ok {
			graph.Tunnels[tunnel.From] = make([]string, 0)
		}
		graph.Tunnels[tunnel.From] = append(graph.Tunnels[tunnel.From], tunnel.To)

		if _, ok := graph.Tunnels[tunnel.To]; !ok {
			graph.Tunnels[tunnel.To] = make([]string, 0)
		}
		graph.Tunnels[tunnel.To] = append(graph.Tunnels[tunnel.To], tunnel.From)

		if _, ok := graph.Capacity[tunnel.From]; !ok {
			graph.Capacity[tunnel.From] = make(map[string]int)
		}
		graph.Capacity[tunnel.From][tunnel.To] = 1

		if _, ok := graph.Capacity[tunnel.To]; !ok {
			graph.Capacity[tunnel.To] = make(map[string]int)
		}
		graph.Capacity[tunnel.To][tunnel.From] = 1
	}

	return graph, nil
}

// String returns a string representation of the Graph
func (g *Graph) String() string {
	result := "\nNodes:\n"
	for name, room := range g.Nodes {
		result += fmt.Sprintf("%s: %+v\n", name, *room)
	}
	result += "\nTunnels:\n"
	for from, toList := range g.Tunnels {
		for _, to := range toList {
			result += fmt.Sprintf("%s -> %s\n", from, to)
		}
	}
	return result
}

// copyGraph creates a deep copy of the capacity graph
func copyGraph(original map[string]map[string]int) map[string]map[string]int {
	copy := make(map[string]map[string]int)
	for u, neighbors := range original {
		copy[u] = make(map[string]int)
		for v, capacity := range neighbors {
			copy[u][v] = capacity
		}
	}
	return copy
}

func (g *Graph) FindAllPaths(start, end string) [][]string {
	var paths [][]string
	residualGraph := copyGraph(g.Capacity)
	visited := make(map[string]bool)
	path := []string{}
	dfs(residualGraph, g.Tunnels, start, end, visited, &path, &paths)
	return paths
}

// FindNonOverlappingPaths finds non-overlapping paths, preferring shorter ones when initial steps overlap
func (g *Graph) FindNonOverlappingPaths(ants int) [][]string {
	allPaths := g.FindAllPaths(g.Start, g.End)
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})
	a := FiltreleYollar(allPaths, ants)
	return a
}

// Yolları filtreler ve çakışan odaları çıkarır
func FiltreleYollar(yollar [][]string, karincaSayisi int) [][]string {
	var filtrelenmisYollar [][]string

	// İki yolun ara odalarda çakışıp çakışmadığını kontrol eden yardımcı fonksiyon
	yollarCakisiyor := func(yol1, yol2 []string) bool {
		kume := make(map[string]bool)
		for _, oda := range yol1[1 : len(yol1)-1] { // Başlangıç ve bitişi hariç tut
			kume[oda] = true
		}
		for _, oda := range yol2[1 : len(yol2)-1] {
			if kume[oda] {
				return true
			}
		}
		return false
	}

	// Çakışmayan yol kombinasyonlarını bulmak için tüm kombinasyonları dene
	var kombinasyonlar func([][]string, int, []int)
	var enIyiKombinasyon []int
	maxYol := 0

	kombinasyonlar = func(yollar [][]string, indeks int, secili []int) {
		if len(secili) > maxYol {
			maxYol = len(secili)
			enIyiKombinasyon = make([]int, len(secili))
			copy(enIyiKombinasyon, secili)
		}

		for i := indeks; i < len(yollar); i++ {
			cakisiyor := false
			for _, s := range secili {
				if yollarCakisiyor(yollar[s], yollar[i]) {
					cakisiyor = true
					break
				}
			}
			if !cakisiyor {
				secili = append(secili, i)
				kombinasyonlar(yollar, i+1, secili)
				secili = secili[:len(secili)-1]
			}
		}
	}

	kombinasyonlar(yollar, 0, []int{})

	for _, indeks := range enIyiKombinasyon {
		filtrelenmisYollar = append(filtrelenmisYollar, yollar[indeks])
		if len(filtrelenmisYollar) == karincaSayisi {
			break
		}
	}

	return filtrelenmisYollar
}

func printPathLevels(paths [][]string, antCount int) {
	antPositions := make([]int, antCount)
	nodeOccupied := make(map[string]bool)
	antSteps := make([]int, antCount)

	// Initialize ant positions and step counts at the beginning
	for i := 0; i < antCount; i++ {
		antPositions[i] = 1
		antSteps[i] = 1
	}

	round := 1
	startNodeConnections := len(paths)

	for {
		allAntsFinished := true
		roundOutput := []string{}

		antsMovingFromStart := 0

		for i := 0; i < antCount; i++ {
			var pathIndex int
			if i == antCount-1 {
				pathIndex = 0 // Son karınca ilk yolu takip eder
			} else {
				pathIndex = i % len(paths)
			}

			if antPositions[i] >= len(paths[pathIndex]) {
				continue
			}

			if antSteps[i] < len(paths[pathIndex]) {
				nextNode := paths[pathIndex][antPositions[i]]

				if antPositions[i] > 0 && antPositions[i]-1 < len(paths[pathIndex]) {
					nodeOccupied[paths[pathIndex][antPositions[i]-1]] = false
				}

				if antPositions[i] == 1 {
					if antsMovingFromStart >= startNodeConnections {
						continue
					}
					antsMovingFromStart++
				}

				if !nodeOccupied[nextNode] || nextNode == paths[pathIndex][len(paths[pathIndex])-1] {
					roundOutput = append(roundOutput, fmt.Sprintf("L%d-%s", i+1, nextNode))
					nodeOccupied[nextNode] = true
					antPositions[i]++
					antSteps[i]++
				}

				if antPositions[i] < len(paths[pathIndex]) {
					allAntsFinished = false
				}
			} else {
				allAntsFinished = false
			}
		}

		if len(roundOutput) > 0 {
			fmt.Println(strings.Join(roundOutput, " "))
		}
		round++

		if allAntsFinished {
			break
		}
	}
}

// dfs is a depth-first search helper function
func dfs(graph map[string]map[string]int, tunnels map[string][]string, current, target string, visited map[string]bool, path *[]string, paths *[][]string) {
	visited[current] = true
	*path = append(*path, current)

	if current == target {
		*paths = append(*paths, append([]string{}, *path...))
	} else {
		for _, neighbor := range tunnels[current] {
			if !visited[neighbor] && graph[current][neighbor] > 0 {
				dfs(graph, tunnels, neighbor, target, visited, path, paths)
			}
		}
	}

	*path = (*path)[:len(*path)-1]
	visited[current] = false
}

func ReadFile2(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	startTime := time.Now()

	filename := os.Args[1]
	antCount, rooms, tunnels, startRoom, endRoom, err := ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	if startRoom.Name == "" || endRoom.Name == "" {
		fmt.Println("ERROR: invalid data format, no start or end room found")
		return
	}

	if len(rooms) == 0 || len(tunnels) == 0 {
		fmt.Println("ERROR: invalid data format, no rooms or tunnels found")
		return
	}

	// Check for other invalid data format conditions
	// For example: duplicated rooms, links to unknown rooms, rooms with invalid coordinates, etc.

	lines, err := ReadFile2(filename)
	if err != nil {
		fmt.Println("ERROR: File not found")
		return
	}

	graph, err := NewGraph(antCount, rooms, tunnels, startRoom, endRoom)
	if err != nil {
		fmt.Println(err)
		return
	}

	paths := graph.FindNonOverlappingPaths(antCount)
	if len(paths) == 0 {
		fmt.Println("ERROR: invalid data format, no path between ##start and ##end")
		return
	}
	fmt.Println(paths)

	// Print file content
	for _, line := range lines {
		fmt.Println(line)
	}

	fmt.Println() // Empty line
	printPathLevels(paths, antCount)
	endTime := time.Now()

	duration := endTime.Sub(startTime)
	if filename == "example06.txt" || filename == "example07.txt" {
		fmt.Printf("Gerçek zamanlı süre: %.10f saniye\n", duration.Seconds())
	}
}

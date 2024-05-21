package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Room struct {
	name  string
	x, y  int
	links []string
}

type Farm struct {
	rooms   map[string]*Room
	start   string
	end     string
	numAnts int
}

func NewFarm() *Farm {
	return &Farm{
		rooms: make(map[string]*Room),
	}
}

func (f *Farm) AddRoom(name string, x, y int) {
	f.rooms[name] = &Room{name: name, x: x, y: y, links: []string{}}
}

func (f *Farm) AddLink(from, to string) {
	f.rooms[from].links = append(f.rooms[from].links, to)
	f.rooms[to].links = append(f.rooms[to].links, from)
}

func (f *Farm) ReadInput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ERROR: Could not open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var phase string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			if line == "##start" {
				phase = "start"
			} else if line == "##end" {
				phase = "end"
			}
			continue
		}

		if phase == "start" {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return fmt.Errorf("ERROR: invalid room format for start room")
			}
			f.AddRoom(parts[0], atoi(parts[1]), atoi(parts[2]))
			f.start = parts[0]
			phase = ""
		} else if phase == "end" {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return fmt.Errorf("ERROR: invalid room format for end room")
			}
			f.AddRoom(parts[0], atoi(parts[1]), atoi(parts[2]))
			f.end = parts[0]
			phase = ""
		} else {
			if strings.Contains(line, "-") {
				parts := strings.Split(line, "-")
				if len(parts) != 2 {
					return fmt.Errorf("ERROR: invalid link format")
				}
				f.AddLink(parts[0], parts[1])
			} else if f.numAnts == 0 {
				f.numAnts = atoi(line)
				if f.numAnts <= 0 {
					return fmt.Errorf("ERROR: No ants specified or invalid number of ants")
				}
			} else {
				parts := strings.Split(line, " ")
				if len(parts) != 3 {
					return fmt.Errorf("ERROR: invalid room format")
				}
				f.AddRoom(parts[0], atoi(parts[1]), atoi(parts[2]))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ERROR: Scanner error")
	}

	if f.start == "" || f.end == "" {
		return fmt.Errorf("ERROR: start or end room not defined")
	}

	return nil
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("ERROR: invalid number format"))
	}
	return n
}

func (f *Farm) DFSAllPaths(start, end string) [][]string {
	var paths [][]string
	var dfs func(currentRoom string, path []string, visited map[string]bool)

	dfs = func(currentRoom string, path []string, visited map[string]bool) {
		path = append(path, currentRoom)
		if currentRoom == end {
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			paths = append(paths, pathCopy)
			return
		}

		visited[currentRoom] = true
		for _, neighbor := range f.rooms[currentRoom].links {
			if !visited[neighbor] {
				dfs(neighbor, path, visited)
			}
		}
		visited[currentRoom] = false
	}

	dfs(start, []string{}, make(map[string]bool))
	return paths
}

func pathLength(path []string) int {
	return len(path)
}

func pathsSimilarity(path1, path2 []string) int {
	similarity := 0
	for i := 0; i < len(path1) && i < len(path2); i++ {
		if path1[i] == path2[i] {
			similarity++
		}
	}
	return similarity
}

func filterPaths(paths [][]string, similarityThreshold int) [][]string {
	var filteredPaths [][]string
	for _, path := range paths {
		isUnique := true
		for _, filteredPath := range filteredPaths {
			if pathsSimilarity(path, filteredPath) >= similarityThreshold {
				isUnique = false
				break
			}
		}
		if isUnique {
			filteredPaths = append(filteredPaths, path)
		}
	}
	return filteredPaths
}

type Ant struct {
	id    int
	path  []string
	index int
}

func assignAntsToPaths(numAnts int, paths [][]string) []Ant {
	ants := make([]Ant, numAnts)
	pathIndex := 0
	for i := 0; i < numAnts; i++ {
		ants[i] = Ant{
			id:    i + 1,
			path:  paths[pathIndex],
			index: 0,
		}
		pathIndex = (pathIndex + 1) % len(paths)
	}
	return ants
}

func printAntMovements(ants []Ant, numAnts int, start, end string) {
	nodeOccupied := make(map[string]bool)

	turn := 0
	for {
		allAntsFinished := true
		roundOutput := []string{}

		for i := 0; i < numAnts; i++ {
			ant := &ants[i]
			if ant.index < len(ant.path)-1 {
				nextNode := ant.path[ant.index+1]
				if !nodeOccupied[nextNode] || nextNode == end {
					if ant.index > 0 {
						nodeOccupied[ant.path[ant.index]] = false
					}
					nodeOccupied[nextNode] = true
					ant.index++
					roundOutput = append(roundOutput, fmt.Sprintf("L%d-%s", ant.id, nextNode))
					allAntsFinished = false
				}
			}
		}

		if len(roundOutput) > 0 {
			fmt.Println(strings.Join(roundOutput, " "))
		}

		if allAntsFinished {
			break
		}
		turn++
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: lem-in <filename>")
		return
	}

	filename := os.Args[1]
	farm := NewFarm()
	err := farm.ReadInput(filename)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	paths := farm.DFSAllPaths(farm.start, farm.end)
	if len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Sort paths by length
	sort.Slice(paths, func(i, j int) bool {
		return pathLength(paths[i]) < pathLength(paths[j])
	})

	// Print the number of ants
	fmt.Println(farm.numAnts)

	// Print the rooms
	fmt.Println("the_rooms")

	fmt.Printf("%s %d %d\n", farm.start, farm.rooms[farm.start].x, farm.rooms[farm.start].y)
	for name, room := range farm.rooms {
		if name != farm.start && name != farm.end {
			fmt.Printf("%s %d %d\n", name, room.x, room.y)
		}
	}
	fmt.Printf("%s %d %d\n", farm.end, farm.rooms[farm.end].x, farm.rooms[farm.end].y)

	// Print the links
	fmt.Println("the_links")

	for from, room := range farm.rooms {
		for _, to := range room.links {
			if from < to {
				fmt.Printf("%s-%s\n", from, to)

			}
		}
	}
	fmt.Println()

	// Filter paths to remove similar ones
	similarityThreshold := 2 // Threshold can be adjusted as needed
	filteredPaths := filterPaths(paths, similarityThreshold)

	// Assign ants to paths
	antPaths := assignAntsToPaths(farm.numAnts, filteredPaths)

	// Print each turn
	printAntMovements(antPaths, farm.numAnts, farm.start, farm.end)
}

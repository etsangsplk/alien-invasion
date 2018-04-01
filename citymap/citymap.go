package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CityNode is a node in our CityMap
type CityNode struct {
	name string
}

func (cn *CityNode) String() string {
	return fmt.Sprintf("%v", cn.name)
}

// CityMap is a graph representation of our city, with
// directions NWSE
type CityMap struct {
	// map: string -> *CityNode
	cities map[string]*CityNode
	// map: CityNode -> map: string (NWSE) -> CityNode
	connections map[CityNode]map[string]*CityNode
	// Let's assume that no one will add more than 4 connections.
	lock sync.RWMutex
}

// AddCity takes an existing CityMap and adds a constructed
// CityNode to the list of cities
func (cm *CityMap) AddCity(city *CityNode) {
	cm.lock.Lock()
	if cm.cities == nil {
		cm.cities = make(map[string]*CityNode)
	}
	name := city.name
	cm.cities[name] = city
	// Instantiate the map if not yet done
	if cm.connections == nil {
		cm.connections = make(map[CityNode]map[string]*CityNode)
	}
	// Instantiate the node's connections
	cm.connections[*city] = map[string]*CityNode{"north": nil, "west": nil, "south": nil, "east": nil}

	cm.lock.Unlock()
}

// AddConnection returns True
func (cm *CityMap) AddConnection(cityname1 string, cityname2 string, direction string) {
	// Let's assume we're not really worried about deadlock by having someone input
	// invalid input and cause some portion to fail in between
	cm.lock.Lock()
	c1 := cm.cities[cityname1]
	c2 := cm.cities[cityname2]

	// Let's assume that anyone who constructs a graph
	// doesn't care if overrides happen
	switch direction {
	case "north":
		cm.connections[*c1]["north"] = c2
		cm.connections[*c2]["south"] = c1
	case "west":
		cm.connections[*c1]["west"] = c2
		cm.connections[*c2]["east"] = c1
	case "south":
		cm.connections[*c1]["south"] = c2
		cm.connections[*c2]["north"] = c1
	case "east":
		cm.connections[*c1]["east"] = c2
		cm.connections[*c2]["west"] = c1
	}
	cm.lock.Unlock()
}

// RemoveCity removes the City and its connections (both ways) in the CityMap
func (cm *CityMap) RemoveCity(cityname string) {
	// Let's assume we're not really worried about deadlock by having someone input
	// invalid input and cause some portion to fail in between
	cm.lock.Lock()
	c1 := cm.cities[cityname]

	// Remove the city from all connections
	c1Connections := cm.connections[*c1]
	for _, c2 := range c1Connections {
		if c2 != nil {
			c2Connections := cm.connections[*c2]
			// Remove c1 from c2's conncetions
			for direction := range c2Connections {
				c2Neighbor := c2Connections[direction]
				// Check if the pointers are the same
				if c2Neighbor == c1 {
					c2Connections[direction] = nil
					break
				}
			}
		}
	}

	// Remove the cities from the list of cities
	delete(cm.cities, cityname)
	delete(cm.connections, *c1)

	cm.lock.Unlock()
}

// PrintMap prints the cities along with their neighbors
func (cm *CityMap) PrintMap() {
	cm.lock.RLock()
	// fmt.Println(len(cm.connections))
	// Sort the keys of cityname -> city mapping
	names := make([]string, 0)
	for c := range cm.cities {
		names = append(names, c)
	}
	sort.Strings(names)

	for _, n := range names {
		city := cm.cities[n]
		connections := cm.connections[*city]

		fmt.Print("CITY: ", city)
		fmt.Print("  CONNECTIONS:")
		for direction, neighborCity := range connections {
			fmt.Printf(" %v=%v", direction, neighborCity)
		}
		fmt.Println()
	}
	fmt.Println()

	cm.lock.RUnlock()
}

// ReadCityMapFile takes in a filename and constructs a citymap from text
func (cm *CityMap) ReadCityMapFile(filename string) *CityMap {
	// We assume that city names can't have spaces
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return cm
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// This is our buffer now
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines {
		cityAndConnections := strings.Split(line, " ")
		// Pull out the cityname and its connections
		c1Name := cityAndConnections[0]
		c1Connections := cityAndConnections[1:]

		// Create the city
		c1 := CityNode{c1Name}

		// Easy add if we're dealing with the first city in the map
		if cm.cities == nil {
			cm.AddCity(&c1)
		} else {
			_, exists := cm.cities[c1Name]
			if !exists {
				cm.AddCity(&c1)
			}
		}

		for _, con := range c1Connections {
			dirAndName := strings.Split(con, "=")
			direction, c2Name := dirAndName[0], dirAndName[1]
			_, exists := cm.cities[c2Name]
			if !exists {
				c2 := CityNode{c2Name}
				cm.AddCity(&c2)
			}
			cm.AddConnection(c1Name, c2Name, direction)
		}
	}

	return cm
}

// PickRandomCity picks a random city from the CityMap
func (cm *CityMap) PickRandomCity() *CityNode {
	cities := make([]*CityNode, len(cm.cities))
	i := 0
	for _, city := range cm.cities {
		cities[i] = city
		i++
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	randCityIdx := r.Intn(len(cities))
	return cities[randCityIdx]
}

// PickRandomNeighbor picks a random node from a mapping of directions
// to cities. Only call if the city you're at has neighbors.
func (cm *CityMap) PickRandomNeighbor(city *CityNode) *CityNode {
	neighborCitiesMap := cm.connections[*city]
	neighborCities := make([]*CityNode, 0)
	for _, city := range neighborCitiesMap {
		if city != nil {
			neighborCities = append(neighborCities, city)
		}
	}
	// fmt.Println(neighborCities)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	randNeighborIdx := r.Intn(len(neighborCities))
	return neighborCities[randNeighborIdx]
}

// hasNeighbors checks if a city has any neighboring cities
func (cm *CityMap) hasNeighbors(city *CityNode) bool {
	neighborCitiesMap := cm.connections[*city]
	if neighborCitiesMap == nil {
		return false
	}

	for _, neighborCity := range neighborCitiesMap {
		if neighborCity != nil {
			return true
		}
	}
	// all of the neighbors were nil
	return false
}

// makeRange takes a min and max and gives us a slice with
// a range of numbers from min to max
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// RunAlienSim runs a simulation of alien invasion for n aliens
// on the CityMap. Assumes that more than two aliens can be in a city
// once
func (cm *CityMap) RunAlienSim(numAliens int) {
	nodesToOccupants := make(map[*CityNode][]int)
	// Assign aliens to random cities
	aliens := makeRange(1, numAliens)
	for _, a := range aliens {
		randomCity := cm.PickRandomCity()
		nodesToOccupants[randomCity] = append(nodesToOccupants[randomCity], a)
		// fmt.Println(nodesToOccupants)
	}

	// Step with aliens
	for steps := 0; steps < 10000; steps++ {
		// Look at all cities with alien occupants
		for city, cityOccupants := range nodesToOccupants {
			// If those cities have neighbors, we can move the occcupants one step
			if cm.hasNeighbors(city) {
				for i := len(cityOccupants) - 1; i >= 0; i-- {
					neighborCity := cm.PickRandomNeighbor(city)
					// fmt.Println(neighborCity)
					// Update the neighboring city's slice of occupants
					// fmt.Println(nodesToOccupants[neighborCity])
					nodesToOccupants[neighborCity] = append(nodesToOccupants[neighborCity], cityOccupants[i])
					// Remove the alien from the present city slice of occupants
					cityOccupants = append(cityOccupants[:i], cityOccupants[i+1:]...)
					// nodesToOccupants[city] = cityOccupants // Update (is this needed?)
				}
			}
		}

		// After the movement has occured for the step
		// we must evaluate the current state and delete any CityNodes
		// with multiple occupants
		for city, cityOccupants := range nodesToOccupants {
			if len(cityOccupants) > 1 {
				fmt.Print(city.name, " has been destroyed by")
				for occupantIdx, occupant := range cityOccupants {
					if occupantIdx == len(cityOccupants)-1 {
						fmt.Printf(" and alien %v!\n", occupant)
					} else {
						fmt.Printf(" alien %v,", occupant)
					}
				}
				cm.RemoveCity(city.name)
				delete(nodesToOccupants, city)
			}
		}
	}

}

func main() {
	// run with ./citymap n map.txt
	num := os.Args[1]
	file := os.Args[2]

	numAliens, _ := strconv.Atoi(num)

	var cm CityMap
	cm.ReadCityMapFile(file)
	cm.RunAlienSim(numAliens)
}

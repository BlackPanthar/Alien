package citymap

//package main

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
	mapset "github.com/deckarep/golang-set/v2"
	
)

// CityNode is a node in CityMap
type cityNode struct {
	name string
}

func (c *cityNode) String() string {
	return fmt.Sprintf("%v", c.name)
}

// CityMap is a graph representation of the city, with North, South, West, East directions
type cityMap struct {
	// the map is string -> *cityNode
	cities map[string]*cityNode
	// map: cityNode -> map: string (North, South, West, East) -> CityNode
	connections map[cityNode]map[string]*cityNode
	// Let's assume that no one will add more than 4 connections.
	lock sync.RWMutex
}

// addCity takes an existing cityMap and adds a constructed
// cityNode to the list of cities
func (c *cityMap) addCity(cn *cityNode) {
	c.lock.Lock()
	if c.cities == nil {
		c.cities = make(map[string]*cityNode)
	}
	name := cn.name
	c.cities[name] = cn
	// if not yet done, Instantiate the map 
	if c.connections == nil {
		c.connections = make(map[cityNode]map[string]*cityNode)
	}
	// Instantiate the nodes connections
	c.connections[*cn] = map[string]*cityNode{"north": nil, "west": nil, "south": nil, "east": nil}

	c.lock.Unlock()
}


func (c *cityMap) addConnection(cityname1 string, cityname2 string, direction string) {
	// Let's assume we're not really worried about deadlock from an invalid input causing some portion to fail
	c.lock.Lock()
	c1 := c.cities[cityname1]
	c2 := c.cities[cityname2]

	// Assuming anyone who constructs a graph doesn't care if overrides happen
	switch direction {
	case "north":
		c.connections[*c1]["north"] = c2
		c.connections[*c2]["south"] = c1
	case "west":
		c.connections[*c1]["west"] = c2
		c.connections[*c2]["east"] = c1
	case "south":
		c.connections[*c1]["south"] = c2
		c.connections[*c2]["north"] = c1
	case "east":
		c.connections[*c1]["east"] = c2
		c.connections[*c2]["west"] = c1
	}
	c.lock.Unlock()
	
	
}

// RemoveCity removes the City and its connections (both ways) in the CityMap
func (c *cityMap) RemoveCity(name string) {
	// Let's assume we're not really worried about deadlock by having someone input
	// invalid input and cause some portion to fail in between
	c.lock.Lock()
	c1 := c.cities[name]

	// Remove the city from all connections
	c1Connections := c.connections[*c1]
	for _, c2 := range c1Connections {
		if c2 != nil {
			c2Connections := c.connections[*c2]
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
	delete(c.cities, name)
	delete(c.connections, *c1)

	c.lock.Unlock()
}

// PrintMap prints the cities along with their neighbors
func (c *cityMap) PrintMap() {
	c.lock.RLock()
	// fmt.Println(len(cm.connections))
	// Sort the keys of cityname -> city mapping
	names := make([]string, 0)
	for c := range c.cities {
		names = append(names, c)
	}
	sort.Strings(names)

	for _, n := range names {
		city := c.cities[n]
		connections := c.connections[*city]

		fmt.Print("CITY: ", city)
		fmt.Print("  CONNECTIONS:")
		for direction, neighborCity := range connections {
			fmt.Printf(" %v=%v", direction, neighborCity)
		}
		fmt.Println()
	}
	fmt.Println()

	c.lock.RUnlock()
}

// ReadCityMapFile takes in a file Name and constructs a citymap from text
func (c *cityMap) ReadCityMapFile(fileName string) *cityMap {
	// We assume that city names can't have spaces
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return c
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// our buffer now
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
		c1 := cityNode{c1Name}

		// Easy add if we're dealing with the first city in the map
		if c.cities == nil {
			c.addCity(&c1)
		} else {
			_, exists := c.cities[c1Name]
			if !exists {
				c.addCity(&c1)
			}
		}

		for _, con := range c1Connections {
			dirAndName := strings.Split(con, "=")
			direction, c2Name := dirAndName[0], dirAndName[1]
			_, exists := c.cities[c2Name]
			if !exists {
				c2 := cityNode{c2Name}
				c.addCity(&c2)
			}
			c.addConnection(c1Name, c2Name, direction)
		}
	}

	return c
}

// PickRandomCity picks a random city from the CityMap
func (c *cityMap) PickRandomCity() *cityNode {
	cities := make([]*cityNode, len(c.cities))
	i := 0
	for _, cityname := range c.cities {
		cities[i] = cityname
		i++
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	randCityIdx := r.Intn(len(cities))
	return cities[randCityIdx]
}

// PickRandomNeighbor picks a random node from a mapping of directions
// to cities. Only call if the city has neighbors.
func (c *cityMap) PickRandomNeighbor(cn *cityNode) *cityNode {
	neighborCitiesMap := c.connections[*cn]
	neighborCities := make([]*cityNode, 0)
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
func (c *cityMap) hasNeighbors(cn *cityNode) bool {
	neighborCitiesMap := c.connections[*cn]
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

// RunAlienSim runs a simulation of alien invasion for N aliens
// on the CityMap. Assumes that more than two aliens can be in a city same time
func (c *cityMap) RunAlienSim(numberOfAliens int) {
	nodesToOccupants := make(map[*cityNode]mapset.Set)
	// Assign aliens to random cities
	aliens := makeRange(1, numberOfAliens)
	for _, a := range aliens {
		randomCity := c.PickRandomCity()
		_, ok := nodesToOccupants[randomCity]
		if ok {
			nodesToOccupants[randomCity].Add(a)
		} else {
			occupantSet := mapset.NewSet()
			occupantSet.Add(a)
			nodesToOccupants[randomCity] = occupantSet
		}
		// fmt.Println(nodesToOccupants)
	}

	
	for steps := 0; steps < 10000; steps++ {
		// Look at all cities with alien occupants
		for city, cityOccupants := range nodesToOccupants {
			// If those cities have neighbors, we can move the occcupants one step
			if c.hasNeighbors(city) {
				it := cityOccupants.Iterator()
				occcupantsToRemove := make([]interface{}, 0)
				for cityOccupant := range it.C {
					neighborCity := c.PickRandomNeighbor(city)
					// fmt.Println(neighborCity)
					// Update the neighboring city's slice of occupants
					// fmt.Println(nodesToOccupants[neighborCity])
					_, ok := nodesToOccupants[neighborCity]
					if ok {
						nodesToOccupants[neighborCity].Add(cityOccupant)
					} else {
						occupantSet := mapset.NewSet()
						occupantSet.Add(cityOccupant)
						nodesToOccupants[neighborCity] = occupantSet
					}
					// Remove the alien from the present city slice of occupants
					occcupantsToRemove = append(occcupantsToRemove, cityOccupant)
				}

				for removedAlien := range occcupantsToRemove {
					cityOccupants.Remove(removedAlien)
				}
			}
		}

		// After the movement has occured for the step/iteration
		// evaluate the current state and delete any CityNodes with multiple occupants
		for city, cityOccupants := range nodesToOccupants {
			if cityOccupants.Cardinality() > 1 {
				fmt.Print(city.name, " has been destroyed by")
				it := cityOccupants.Iterator()

				count := 0
				for cityOccupant := range it.C {
					if count == cityOccupants.Cardinality()-1 {
						fmt.Printf(" and alien %v!\n", cityOccupant)
					} else {
						fmt.Printf(" alien %v,", cityOccupant)
					}
					count++
				}
				c.RemoveCity(city.name)
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

	var c cityMap
	c.ReadCityMapFile(file)
	c.RunAlienSim(numAliens)
}

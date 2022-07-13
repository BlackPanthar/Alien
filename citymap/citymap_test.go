package main

import (
	"fmt"
	"testing"
)

func fillCityAndRemove() {
	var c cityMap
	nA := cityNode{"A"}
	nB := cityNode{"B"}
	nC := cityNode{"C"}
	nD := cityNode{"D"}
	nE := cityNode{"E"}
	nF := cityNode{"F"}
	c.addCity(&nA)
	c.addCity(&nB)
	c.addCity(&nC)
	c.addCity(&nD)
	c.addCity(&nE)
	c.addCity(&nF)

	c.addConnection("A", "B", "west")
	c.addConnection("A", "C", "south")
	c.addConnection("B", "E", "west")
	c.addConnection("A", "D", "north")
	fmt.Println("===Map After Connections Added===")
	c.PrintMap()

	c.RemoveCity("A")
	fmt.Println("===Map After City \"A\" Removed===")
	c.PrintMap()

}

func TestAdd(t *testing.T) {
	// fillCityAndRemove()
	var c cityMap
	c.ReadCityMapFile("map.txt")
	c.RunAlienSim(2)
	fmt.Println("\n===Map after Alien Sim ===")

}

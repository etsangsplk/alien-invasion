package main

import (
	"fmt"
	"testing"
)

func fillCityAndRemove() {
	var cm CityMap
	nA := CityNode{"A"}
	nB := CityNode{"B"}
	nC := CityNode{"C"}
	nD := CityNode{"D"}
	nE := CityNode{"E"}
	nF := CityNode{"F"}
	cm.AddCity(&nA)
	cm.AddCity(&nB)
	cm.AddCity(&nC)
	cm.AddCity(&nD)
	cm.AddCity(&nE)
	cm.AddCity(&nF)

	cm.AddConnection("A", "B", "west")
	cm.AddConnection("A", "C", "south")
	cm.AddConnection("B", "E", "west")
	cm.AddConnection("A", "D", "north")
	fmt.Println("===Map After Connections Added===")
	cm.PrintMap()

	cm.RemoveCity("A")
	// cm.RemoveCity("B")
	// cm.RemoveCity("C")
	// cm.RemoveCity("D")
	fmt.Println("===Map After City \"A\" Removed===")
	cm.PrintMap()

}

func TestAdd(t *testing.T) {
	// fillCityAndRemove()
	var cm CityMap
	cm.ReadCityMapFile("map.txt")
	fmt.Println("===Map From Map.txt Added===")
	cm.PrintMap()
	cm.RunAlienSim(2)
	fmt.Println("\n===Map after Alien Sim ===")
	cm.PrintMap()

}

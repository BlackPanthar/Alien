# Alien


# Alien Invasion 

I immediately thought about using graphs. We can have each node represent a city that aliens can visit.

So started off with a very simple graph implementation in Go. Made sure it was thread-safe
by using locks

### The CityMap

I used a map to keep track of the cities and connected cities. I used a mapping with the strings of each direction (north, west, south, east) to whichever CityNode was located in that direction. As a result, I had a CityMap that now had CityNodes connected via direction. 


### Alien Invasion

For the alien invasion, I kept it simple by adding functions inside of the cityMap.go file. For this, I assumed that more than two aliens can be in a city at once. I created a mapping that would take a cityNode and provide a list (slice) of alien occupants. From there I could keep track of which aliens were in which cities, manage those aliens every step.

I would then assign aliens to random cities, and aliens were just integers in the list of occupants, so it was easy to manage.

After that I started the simulation with a limit of 10000 steps, and I would look at the cities with occupants and then check if they had neighbors. If they had neighbors, then we'd pick a random one, update the neighbor's occupant list and then update the current city's occupant list.

After movements had occurred for the step, then we needed to evaluate what happened at the step. I would go through all CityNodes with occupants and check if their occupant list exceeded 1. If it did, that would result in a battle between aliens and the destruction of the city. 

I couldn't really test this portion formally because it relied on randomness, but in production I would remove the aspect of randomness to make sure the rest of the code worked, and then I would be able to test a non-random version of the alien simulation.

I realized a bug where I would be getting duplicate elements per CityNode in terms of occupants.

I solved this by using a Set data structure here: https://github.com/deckarep/golang-set

## How to Run

Command line program in citymap directory:
`go run citymap.go numberOfAliens filename`

```

➜  citymap git:(main) ✗ go run citymap.go 6 map.txt
Bee} has been destroyed by alien 5, alien 1, and alien 4!
Foo has been destroyed by alien 3, alien 2, and alien 6!


➜  citymap git:(main) ✗ go test
Baz has been destroyed by alien 1, and alien 2!

===Map after Alien Sim ===
PASS
ok      Alien/citymap   0.671s

➜  graph git:(main) ✗ go test
PASS
ok      Alien/graph     0.293s
```


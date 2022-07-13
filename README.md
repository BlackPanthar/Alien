# Alien


# Alien Invasion 

## Thought Process

My immediate reaction when reading the spec was the idea of a graph, because that's what
all city maps really are. We can have each node represent a city that aliens can visit.

So started off with a very simple graph implementation in Go. Made sure it was thread-safe
by using locks

### The CityMap

I used a map to keep track of the cities and connected cities. I used a mapping with the strings of each direction (north, west, south, east) to whichever CityNode was located in that direction. As a result, I had a CityMap that now had CityNodes connected via direction. 


### Alien Invasion

For the alien invasion, I kept it simple by adding functions inside of the CityMap.go file. For this, I assumed that more than two aliens can be in a city at once. I created a mapping that would take a CityNode and provide a list (slice) of alien occupants. From there I could keep track of which aliens were in which cities, manage those aliens every step.

I would then assign aliens to random cities, and aliens were just integers in the list of occupants, so it was easy to manage.

After that I started the simulation with a limit of 10000 steps, and I would look at the cities with occupants and then check if they had neighbors. If they had neighbors, then we'd pick a random one, update the neighbor's occupant list and then update the current city's occupant list.

After movements had occurred for the step, then we needed to evaluate what happened at the step. I would go through all CityNodes with occupants and check if their occupant list exceeded 1. If it did, that would result in a battle between aliens and the destruction of the city. 

I couldn't really test this portion formally because it relied on randomness, but in production I would remove the aspect of randomness to make sure the rest of the code worked, and then I would be able to test a non-random version of the alien simulation.

I realized a bug where I would be getting duplicate elements per CityNode in terms of occupants.
Here's an example of something that doesn't make sense.
```
❯ ./citymap 4 map.txt
Bee has been destroyed by alien 1, alien 2, and alien 4!
Bar has been destroyed by alien 4, alien 2, and alien 1!
Foo has been destroyed by alien 3, and alien 3!
Baz has been destroyed by alien 3, alien 3, and alien 3!
```

I solved this by using a Set data structure here: https://github.com/deckarep/golang-set/v2

## How to Run

Command line program in citymap directory:
`./citymap numberOfAliens filename`

```
Alien/citymap git/main
❯ ./citymap 5 map.txt
Bar has been destroyed by alien 3, alien 5, and alien 1!
Qu-ux has been destroyed by alien 2, alien 4, and alien 5!
Foo has been destroyed by alien 5, and alien 4!
```


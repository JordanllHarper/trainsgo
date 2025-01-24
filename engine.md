# Engine layer

This layer will comprise of the simulation aspect for Trainsgo. Keeping track of trains coordinates, updating the database accordingly and updating train states such as when they're at a location and loading or unloading passengers.

## Requirements

I want to:
- set a train from one station to go to another station.
- view the status as it travels in real time

Assumptions:
- there will be no "lines". Trains can travel in the most efficient manner (diagonally, straight line). This would be a good change after we have a prototype working.
- There will be constant loading and unloading times. 
- No "journeys". Point a to point b only. This is another good change.


Problems:
- how are we representing a train? A grid data structure?
- how are we pathfinding a train to a station?
- how are we taking speed into account?
- how are we keeping track of trains moving?
- how are we updating the simulation in real time? Perhaps channels can be used here?

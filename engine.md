# Engine layer

This layer will comprise of the simulation aspect for Trainsgo. Keeping track of trains coordinates, updating the database accordingly and updating train states such as when they're at a location and loading or unloading passengers.

## Requirements for prototype

- set a train from one station to go to another station.
- view the status as it travels in real time

## Assumptions

- there will be no "lines". Trains can travel in the most efficient manner (diagonally, straight line). This would be a good change after we have a prototype working.
- There will be constant loading and unloading times. 
- No "journeys". Point a to point b only. This is another good change.


## Problems

- how are we representing a train? A grid data structure?
- how are we pathfinding a train to a station?
- how are we taking speed into account?
- how are we keeping track of trains moving in real time? Are we updating the database constantly? Is this perhaps not a good idea?
- how are we updating the simulation in real time? Perhaps channels can be used here? How are we sending events (such as starting trains on journeys)?
- how are we making sure the simulation follows the time correctly? i.e. we can't just have a for loop, as that will cause trains to update as fast as the for loop. We might need to pad runtime to make sure things are kept in sync with timing. Not sure on this issue. 


## Thoughts

High level stuff: 

- we need a channel between our backend and the engine, through which we can send events through.
- a user will send a post request for a new trip to start. They will define the starting destination and the ending destination. We will log this in another db table.
- the backend service will pick this request up and validate it is a good request - the train is unused in the simplest case. Then it will send along a channel into the simulation.
- a goroutine will pick that up and start the journey from the designated station to the target station. The train entity will be created and travel at it's appropriate speed. We'll say every "cell" of the grid is 1km. 1000m in 1km, so train travelling at 500 m/s would get to the destination in 2 seconds.
- trains will report their current coordinates to some sort of receiver - not sure on details here. 
- another goroutine will report the "state" of the simulation from above receiver. 
- we can then send the state of the simulation back out - not sure how we would do real time request responses here though...


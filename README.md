### Business Case

- Send driver's location updates from mobile and forward them to a broker (let's say every 10s)
- Store locations, compute if a driver is moving or not
- A driver status is not moving if he has driven less than 500 meters in the last 5 minutes
    // sliding window algorithm



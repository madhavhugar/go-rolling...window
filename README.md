## Description

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.


### TODOs

- ensure we write to a file with some mutex/semaphone support
- ensure that we have a inbuilt clock which triggers every minute to reset the counter

### Other TODOs:

- ensure we handle errors nicely

cannot handle inconsistent state

array that stores timestamps
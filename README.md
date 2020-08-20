## Description

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

## Solution walkthrough

To implement a rolling window we take a lazy approach of updating our window only when a request comes in.

On starting the server =>
  - load the log file and populate the window with relevant records

On every incoming request =>
  - garbage collect older records from window
  - append the new record to a log file
  - add the new record to a log file

## Setup


1.  Build the project

```go
go build -race -o cardinality
```

2. Execute the executable

```bash
./cardinality
```

## Enhancements

- better error handling
- handle inconsistent states/avoid assumptions
- separate modules into different packages
- binary search the log file for improved server load time

# Distributed MapReduce System in Go

## Overview
This project implements a simplified distributed MapReduce system in Go, with a web dashboard for real-time system monitoring. It is minimalist, self-contained, and uses only Go's standard library.

---

## Project Structure

```
.
├── common/ # Shared types and utilities
├── mapreduce/ # Mapping and reducing logic
├── master/ # RPC server, task assignment, dashboard
├── worker/ # Worker-side processing logic
├── dashboard/ # Web interface (HTML/CSS/JS)
├── tests/ # Unit and integration tests
├── main.go # Program entry point
├── go.mod # Go module definition
```

---

## Getting Started

### 1. Prerequisites
- Go 1.17 or later
- Terminal & web browser

### 2. Install Dependencies
No external dependencies are required.  
Make sure you have a valid `go.mod` file like this:
```go
module project

### 3. Prepare Input Files
Add some `.txt` files to the root, for example:
```
file1.txt
file2.txt
file3.txt
```

---

##  Execution

###  Start the Master
```bash
go run main.go -mode=master -files="file1.txt,file2.txt,file3.txt"
```
Open the dashboard at : [http://localhost:8080](http://localhost:8080)

###  Start One or More Workers
In one or more other terminals:
```bash
go run main.go -mode=worker
```

---

##  Testing
Run the tests with:
```bash
go test ./tests
```

---

##  Dashboard
- **Web interface**: Simple HTML/CSS/JS for monitoring.
- List of **workers** (address, status)
- List of **tasks** (type, file, status, assigned worker)
- **Dynamic progress bar**
- Top 10 words

---
## Other Features
Along the web interface, you can take advantage of the following features using the URL bar:
- **localhost/8080/results**: View the results of the MapReduce operation.
- **localhost/8080/data**: View the status of tasks.



##  Technical Notes
- Uses `net/rpc` and `net/http` for communication.
- Workers can simulate random crashes/slowdowns.
- No external framework, everything is based on the standard library.

---

##  Authors
- Project carried out as part of the distributed applications module (1st year of engineering degree)
Faculty of Sciences of Bizerte — Academic year 2024/2025

Arbi Aya


---

##  Learning Objectives
Understand the concepts of:
- Concurrent programming
- Distributed architecture (master/worker)
- Web monitoring
- Resilience to failures

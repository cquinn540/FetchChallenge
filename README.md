# Install

Note: This project was developed with Go version 1.24

```bash
# Get the code from github
git clone https://github.com/cquinn540/FetchChallenge.git

# Change to the root directory
cd FetchChallenge

# Build the go binary
go build
```

# Run

In the `FetchChallenge` directory:

```bash
# Run with default hostname and port: 'localhost:8080'
./FetchChallenge

# Or run with a custom hostname and port
./FetchChallenge -h '127.0.0.1' -p 8000
```

# Run Integration Tests

In the `FetchChallenge` directory:

```bash
go test ./router
```

# Installation 
To install dependencies without using $GOROOT, in the top level directory, run
```
go mod vendor
```
To install dependencies to $GOROOT, in the top level directory, run
```
go mod tidy
```
# Build and Run
To build and run, in the top level directory run
```
go run .
```
# Test
To test, in the top level directory run
```
go test
```
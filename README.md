This branch uses goroutines to request data synchronously and can cause rate limit errors. Try the master-no-concurrency branch for syncronous data fetching. 

# Installation 
To install dependencies without using $GOROOT, in the top level directory, run:
```
go mod vendor
```
To install dependencies to $GOROOT, in the top level directory, run:
```
go mod tidy
```

# Configure Enviorment variables
To hit the api you must create a file named .env with your api key. An example of this file can be found in the example.env file. It will look like this:
```
API_KEY=<Your Api Key Here>
```

# Build and Run
To build and run, in the top level directory run:
```
go run .
```
# Test
To test, in the top level directory run:
```
go test
```

Created By Ryan Callahan

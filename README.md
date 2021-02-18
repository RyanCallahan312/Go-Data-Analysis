This branch uses goroutines to request data asynchronously and can cause rate limit errors. Try the master-no-concurrency branch for syncronous data fetching. 

# Installation 
This assumes that Go 1.15 or later and PostgreSQL 13.0 or later are already installed. If they are not installed you should do that first following these guides:
- Install Go: https://golang.org/doc/install
- Install PostgreSQL: https://www.postgresql.org/download/
To install dependencies without using $GOROOT, in the top level directory, run:
```
go mod vendor
```
To install dependencies to $GOROOT, in the top level directory, run:
```
go mod tidy
```
To create the PostgreSQL user that will create the db, in the top level directory, run:
```
sh ./db_scripts/CREATE_INITIAL_USER.sh
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

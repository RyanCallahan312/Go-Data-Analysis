Created By Ryan Callahan

# Installation

This project assumes that Go 1.15 or later and PostgreSQL 13.0 or later and make are already installed. This project also assumes that you are using some kind of bash-like terminal (this is built into macOS and most linux distros). If they are not installed you should do that first following these guides:

-   Install Go: https://golang.org/doc/install
-   Install PostgreSQL: https://www.postgresql.org/download/
-   Install Bash on Windows: https://itsfoss.com/install-bash-on-windows/
-   Install make on Ubuntu/Wsl2 Ubtunu by running: apt-get install build-essential

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

```API_KEY=<Your Api Key Here>
MAINTENANCE_CONNECTION_STRING=port=<Postgres Port Number> host=<Host IP> user=root password=toor<default user created by scripts, you may want to change this> dbname=posgres sslmode=disable
WORKING_CONNECTION_STRING=port=<Postgres Port Number> host=<Host IP> user=root password=toor<default user created by scripts, you may want to change this> dbname=<whatever is specified in the DATABASE_NAME variable>
TEST_CONNECTION_STRING=<same as WORKING_CONNECTION_STRING except the database field should have 'test' added to the end>
DATABASE_NAME=<whatever db name you want, make sure this aligns with your connection strings>
```

***Make sure dbname is the last field in your connection string***

# Build and Run

To build and run, in the top level directory run:

```
make run
```

# Test

To test, in the top level directory run:

```
make test
```

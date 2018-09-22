# README

## SPECIFICATION
For the specification read the SPECIFICATION.md

## BUILDING AND INSTALLING THE SOFTWARE

### Prerequisites

* Golang version 1.10 or later
* Postgres version 9.6 or later
* Docker version 18.06.0-ce or later
* Docker compose version 1.22.0 or later

### Build to binary

Run `make build` and you will have a binary named `go-jwt-login-example` in `out` directory.

### Running the application

After building to binary executable file, you can run just with `./out/go-jwt-login-example`. To make a configuration you can use this program argument:

- `-secret-key` This is a secret key for your application. Example `-secret-key abcdefghijklmnopqrstuvwxyz`
- `-listen-address` Where this application should run. Example `-listen-address localhost:8000`, so you can access the application at localhost:8000
- `-db-url` PostgreSQL connection string. Example `-db-url postgres://username:password@localhost:5432/go-users?sslmode=disable`
- `-db-debug` Whether to show the SQL in log output or not. Default is false. Example `-db-debug true`

### Unit Test

To run unit test, just run `make test` it will also run coverage test.

## REST API documentation

After you running the application binary as mentioned above, you can see the REST API documentation in root path. 
That said, if you bind this application to [localhost:8000](localhost:8000), when you access it will show you an API documentation. There you can try to consume the API right from the documentation and you will see the response in there too.

## Go Documentation

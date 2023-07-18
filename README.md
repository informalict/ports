# About the service

This service allows to create, update or get port's data. 

It starts with some predefined data located in `./assets/ports.jon`.

# Requirements

Tools which should be installed locally:
- make
- docker

# Start service

Run service on 8080 (by default) host port:
```shell
make run
```

Run service on custom host port:
```shell
API_PORT=8081 make run
```

# Exemplary operations on port's service

Get `test` port ID: 
```shell
curl http://localhost:8080/api/v1/ports/test
```

Create `test` port ID:
```shell
curl -X POST --data '{ "name": "test", "country":"test", "coordinates": [1,1] }' http://localhost:8080/api/v1/ports/test
```

Update `test` port ID:
```shell
curl -X PUT --data '{ "name": "new_test", "country":"test", "coordinates": [1,1] }' http://localhost:8080/api/v1/ports/test
```

# For developers

### Start working with this project

Install third-party tools:
```shell
make tools
```

### Testing

Run unit tests:
```shell
make run-unit-tests
```

Run integration tests:
```shell
make run
make run-tests
make clean
```

Check quality of code:
```shell
make linter
```

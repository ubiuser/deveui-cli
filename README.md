## deveui-cli

<p>This is a command line application written in Golang for generation of unique 16-character (hex) identifiers called DevEUI which are used for each MachineMax senor</p>


## How to run this client library

To run this code, first run `touch .env` and add these vars:

```
BASE_URL=http://europe-west1-machinemax-dev-d524.cloudfunctions.net
TIMEOUT=30000
CODE_REGISTRATION_LIMIT=100
```

Then to run locally, use: `go run main.go`.


Alternatively, this code can also be run via [Docker](https://www.docker.com/). To build the docker image use: `docker build -t deveui-cli . --build-arg BASE_URL=${BASE_URL} --build-arg TIMEOUT=${TIMEOUT} --build-arg CODE_REGISTRATION_LIMIT=${CODE_REGISTRATION_LIMIT}` (After touch .env and copying env variables from above).

To run the docker image: `docker run deveui-cli`

To run the tests use: `go test ./...` and to check code coverage: `go test -coverprofile=coverage.out ./... ;    go tool cover -html=coverage.out`

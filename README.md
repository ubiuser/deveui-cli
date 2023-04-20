## deveui-cli

<p>This is a command line application written in Golang for generation of unique 16-character (hex) identifiers called DevEUI.</p>

<p>These DevEUIs are used for each MachineMax senor.</p>


## How to run this client library



To run this code locally, run `touch .env` and add these vars:

```
BASE_URL=http://europe-west1-machinemax-dev-d524.cloudfunctions.net
TIMEOUT=30000
CODE_REGISTRATION_LIMIT=100
```

Then run `go run main.go` to register devices.


Alternatively, this code can also be run via [Docker](https://www.docker.com/). To build the docker image use: `docker build -t deveui-cli . ` 

To run the docker image: `docker run deveui-cli`

Finally, to run the tests use: `go test ./...` and to check code coverage: `go test -coverprofile=coverage.out ./... ;    go tool cover -html=coverage.out`

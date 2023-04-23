<br />
<div align="center">
  <a href="https://machinemax.com/">
    <img src="images/logo.jpeg" alt="MachineMax Logo" width="200" height="200">
  </a>

  <h1 align="center">DevEUI CLI</h1>

  <p align="center">
    <h2 align="center">A Golang program for concurrently registering DevEUI identifiers for MachineMax.</h2>
  </p>
</div>

## About MachineMax DevEUI

Each MachineMax sensor has a unique 16-character (hex) identifier called a DevEUI. As part of
the manufacturing process, it is written onto the internal storage of the sensor. The DevEUI is
also printed on a label on the side of the sensor alongside a 5-character code (the last 5
characters of the DevEUI). For example, a DevEUI of 78111FFFE452555B would have a short
code of 2555B.

The sensors communicate with the MachineMax cloud though a LoRaWAN provider and the
LoRaWAN provider uses the DevEUI to identify the sensor. This means we first have to register
the DevEUI with the provider before we can use it. We pay for every device registered with the
LoRaWAN provider, so it is important that we only register DevEUIs that we use.

When a customer registers a new sensor, they will enter the 5-character short-form code instead
of the full DevEUI, so it is essential that each DevEUI in the batch has a unique 5-char code (for
lookups).

### Built With

* Golang
* Docker

## Getting Started

To run this code, we first need a `.env` file in the root of the project. Once this is done, add these vars:

```
BASE_URL=http://europe-west1-machinemax-dev-d524.cloudfunctions.net
TIMEOUT=30000
CODE_REGISTRATION_LIMIT=100
```

### Prerequisites

- You will need Golang to run this program which can be downloaded at: [https://go.dev/](https://go.dev/)
- This program can also be run using docker, this can be downloaded at: [https://www.docker.com/](https://www.docker.com/)

## Usage

Then to run locally, use: `go run main.go`.

Alternatively, this code can also be run via docker. To build the docker image use: 
```
docker build -t deveui-cli . --build-arg BASE_URL=${BASE_URL} --build-arg TIMEOUT=${TIMEOUT} --build-arg CODE_REGISTRATION_LIMIT=${CODE_REGISTRATION_LIMIT}.
```

To run the docker image: 
```
docker run deveui-cli
```

To run the tests use: 
```
go test ./...
``` 
To check code coverage: 
```
go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out
```

And to run benchmark tests for CPU and memory consumption:
```
go test -bench=. -benchmem
```

Finally to check for race conditions, use: 
```
go test -race ./...
```

## Checklist

- [x] Implement solution
- [x] Write unit tests
- [x] Write readme
- [ ] More unit tests around unhappy paths for higher code coverage

## Contact

Email - nickgowdy87@gmail.com

Website <a href="http://www.nickgowdy.com/" target="_blank">http://www.nickgowdy.com/</a>

Github <a href="https://github.com/nickgowdy" target="_blank">https://github.com/nickgowdy</a>





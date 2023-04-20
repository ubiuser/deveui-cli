FROM golang:1.19-alpine

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/deveui-cli

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./deveui-cli .

# This container exposes port 8080 to the outside world
EXPOSE 8081

# Run the binary program produced by `go install`
ENTRYPOINT ["./deveui-cli"]
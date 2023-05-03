# use the latest Go version and also give it a name so that we can leverage multi-stage builds
FROM golang:1.20-alpine AS base

# why is git added, it doesn't seem to be used
RUN apk add --no-cache git

# it's common practice to create and use a non-root user in the image
RUN adduser -D -g '' nonroot_user

# Set the Current Working Directory inside the container
WORKDIR /app/deveui-cli

# We want to populate the module cache based on the go.{mod,sum} files.
# this won't copy go.sum, just use COPY go.*  .
COPY go.mod .

RUN go mod download

COPY . .

# At this point everything from the golang alpine image will be part of your application, so the image
# size will be fairly large. Instead, add a build step and at the end use scratch and copy the compiled
# files to the final image.
FROM base AS go-builder

# Build the Go app
RUN go build -o ./deveui-cli .

# this will make your final image size considerably smaller
FROM scratch AS production
COPY --from=go-builder ./deveui-cli /deveui-cli

# if you used a non-root user, then add them here
USER nonroot_user

# Run the binary program produced by `go install`
ENTRYPOINT ["/deveui-cli"]

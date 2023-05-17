FROM golang:1.20-alpine as base
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./

FROM base AS go-builder
ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
RUN go build -ldflags "-w -s" -o "./deveui-cli" main.go

FROM scratch AS production
COPY --from=go-builder /app/deveui-cli /bin/
ENTRYPOINT ["/bin/deveui-cli"]

FROM tiltdev/entr AS helper

FROM golang:alpine

COPY --from=helper /tilt-restart-wrapper /usr/local/bin/entr /
RUN date > /.restart-process

WORKDIR /app
COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY *.go ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -o demo-app .

ENTRYPOINT ["/app/demo-app"]

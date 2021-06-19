FROM tiltdev/entr AS helper

FROM golang:alpine

COPY --from=helper /tilt-restart-wrapper /usr/local/bin/entr /
RUN date > /.restart-process

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY *.go ./
RUN GOOS=linux GOARCH=amd64 go build -o demo-app .
ENTRYPOINT ["/app/demo-app"]

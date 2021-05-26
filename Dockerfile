FROM golang:1.16 as builder

WORKDIR /app

# Populate module cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/demo-app .

FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/build/demo-app /app/build/
ENTRYPOINT ["/app/build/demo-app"]
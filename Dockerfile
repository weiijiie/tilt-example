FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY build ./
ENTRYPOINT ["/app/demo-app"]

FROM golang:1.11-alpine3.8 as builder

RUN mkdir -p /go/src/github.com/bladedancer/authz

WORKDIR /go/src/github.com/bladedancer/authz

# Copy necessary files
ADD . .

RUN rm -rf bin
RUN CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o bin/authz ./cmd/authz

# Create non-root user
RUN addgroup -S bladedancer && adduser -S bladedancer -G bladedancer
RUN chown -R bladedancer:bladedancer /go/src/github.com/bladedancer/authz/bin/authz
USER bladedancer

# Base image
FROM scratch

# Copy binary and user from previous build step
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/bladedancer/authz/bin/authz /root/authz
COPY --from=builder /etc/passwd /etc/passwd
USER bladedancer

ENTRYPOINT ["/root/authz"]

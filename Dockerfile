FROM golang:1.11-alpine3.8 as builder

RUN mkdir -p /go/src/github.com/bladedancer/envoyxds

WORKDIR /go/src/github.com/bladedancer/envoyxds

# Copy necessary files
ADD . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/envoyxds ./cmd/envoyxds

# Create non-root user
RUN addgroup -S bladedancer && adduser -S bladedancer -G bladedancer
RUN chown -R bladedancer:bladedancer /go/src/github.com/bladedancer/envoyxds/bin/envoyxds
USER bladedancer

# Base image
FROM scratch

# Copy binary and user from previous build step
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/bladedancer/envoyxds/bin/envoyxds /root/envoyxds
COPY --from=builder /go/src/github.com/bladedancer/envoyxds/envoy /
COPY --from=builder /etc/passwd /etc/passwd
USER bladedancer

ENTRYPOINT ["/root/envoyxds"]

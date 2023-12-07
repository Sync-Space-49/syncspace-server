FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . ./
# This is where one could build the application code as well.
RUN go build -o /main


FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Copy binary to production image.
COPY --from=builder /app/* /app/
COPY --from=builder /main /app/main


# Copy Tailscale binaries from the tailscale image on Docker Hub.
COPY --from=docker.io/tailscale/tailscale:stable /usr/local/bin/tailscaled /app/tailscaled
COPY --from=docker.io/tailscale/tailscale:stable /usr/local/bin/tailscale /app/tailscale
RUN mkdir -p /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

# Run on container startup.
CMD ["sh", "/app/start.sh"]
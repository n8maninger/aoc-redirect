# Use the official Go image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code into the container
COPY . .

# Build the Go application
RUN go generate ./...
RUN go build -o bin/ -tags='netgo timetzdata' -trimpath -a -ldflags '-s -w ' ./cmd/redirect

FROM docker.io/chromedp/headless-shell

COPY --from=builder /app/bin/* /usr/bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


# Run the application
ENTRYPOINT [ "redirect" ]

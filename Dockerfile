FROM golang:1.25 AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
# CGO_ENABLED=0 is critical for creating a statically linked binary
# This binary will run in the very minimal 'scratch' base image.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/app .

# --- STAGE 2: PRODUCTION ---
# Use 'scratch' for the smallest, most secure final image.
# It contains nothing but the kernel and the Go binary.
FROM scratch

# Set environment variable (optional, but good practice)
ENV CGO_ENABLED=0

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/app /usr/local/bin/app
# dont forget the pub key
COPY --from=builder /usr/src/app/authorised /authorised


# Command to run the application
CMD ["/usr/local/bin/app"]

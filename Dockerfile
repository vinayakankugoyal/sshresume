FROM --platform=$BUILDPLATFORM golang:1.25 AS builder

# Set the working directory
WORKDIR /app

# Install dependencies (utilizing Docker cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Use ARG to get target OS and Architecture from Buildx
ARG TARGETOS
ARG TARGETARCH

# Build the binary with cross-compilation flags
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /sshresume ./cmd/sshresume/main.go

# Final stage: a tiny runtime image
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /sshresume /app/sshresume
COPY --from=builder /app/resume /app/resume
ENTRYPOINT ["/app/sshresume"]
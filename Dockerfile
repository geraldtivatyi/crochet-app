# ==========================================
# STAGE 1: Build the Go binary
# ==========================================
FROM golang:1.27-alpine AS builder

# Set the working directory inside the build container
WORKDIR /app

# Copy dependency definition files first
COPY go.mod go.sum* ./

# Download dependencies (cached if go.mod/go.sum haven't changed)
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the binary statically
# CGO_ENABLED=0 creates a standalone binary without dependencies on C libraries
# -ldflags="-s -w" strips debugging info to shrink binary size further
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o crochet-app .

# ==========================================
# STAGE 2: Create the lean production image
# ==========================================
FROM alpine:latest

# Install ca-certificates to allow the app to make outbound HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root group and user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/crochet-app .

# Assign ownership of the application directory to our non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose the default port (documentation step for Docker runtime)
EXPOSE 8080

# Define default environment variables
ENV PORT=8080

# Command to execute when the container starts
CMD ["./crochet-app"]

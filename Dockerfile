FROM golang:1.22.3-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Golang image size reduction optimization 

# Stage 1: Build the Go binary
#FROM golang:1.23.0-bookworm AS build

#WORKDIR /app

# Download Go modules
#COPY go.mod go.sum ./
#RUN go mod download

# Copy all source files and build the binary
#COPY . .

# Build the Go binary with optimizations
#RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
#    -o server -a -ldflags="-s -w" -installsuffix cgo

# Compress the binary using UPX
#RUN upx --ultra-brute -qq server && upx -t server

# Stage 2: Use a minimal base image (scratch) for the final container
#FROM scratch

# Copy the compiled server from the build stage
#COPY --from=build /app/server /server

# Set the entry point to run the Go binary
#ENTRYPOINT ["/server"]


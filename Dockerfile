FROM golang:1.22.3-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

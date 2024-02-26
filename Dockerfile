
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -o rinha ./src/main.go

FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage into the final image
COPY --from=builder /app/rinha .

# Command to run the Go application
EXPOSE 8080

# RUN chmod +x rinha
CMD ["./rinha"]

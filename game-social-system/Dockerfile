# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.21.3 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

# Stage 2: Run
FROM golang:1.21.3
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]

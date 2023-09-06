# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18 
WORKDIR /app
COPY --from=builder /app/main .
# We must copy file with config values to the final container because Viper
# read these values at runtime not during compile in build stage
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]
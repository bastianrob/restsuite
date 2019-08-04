FROM golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o restsuite cmd/*.go


######## Start a new stage from scratch #######
FROM alpine:latest
RUN apk --no-cache add ca-certificates && update-ca-certificates
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/restsuite .

# Expose port 7001 to the outside world
EXPOSE 7001

# Command to run the executable
CMD ["./restsuite"] 
FROM golang:alpine as builder
ENV GO111MODULE=on
# Adding git, bash and openssh to the image
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder ./main .      

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]
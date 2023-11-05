# # Use an official Golang runtime as a parent image
# FROM golang:1.17

# # Set the working directory inside the container
# WORKDIR ./

# # Copy the local package files to the container's workspace
# COPY . .

# # Build the Go application
# RUN go get -d -v ./...
# RUN go install -v ./...

# # Run the Go application
# CMD ["app"]
FROM mysql:8.0
ENV LANG ja_JP.UTF-8
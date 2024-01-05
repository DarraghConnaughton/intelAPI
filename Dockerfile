# Builder Stage
FROM golang:latest AS builder

# Set the working directory in the builder image
WORKDIR /cmd

# Copy the Go source code and Makefile
COPY . .

# Build the Golang binary
RUN make build

# Move executable to final image.
FROM ubuntu:latest

# Set the working directory in the final image
WORKDIR /cmd
# Copy the binary from the builder image to the final image
COPY --from=builder /cmd/releases/intelagent /cmd/
#COPY --from=builder /cmd/data/ /cmd/

# Create a non-root user and set permissions
RUN groupadd -r intelagent && useradd -r -g intelagent intelagent
RUN chown -R intelagent:intelagent /cmd

# Expose port 8080
EXPOSE 8080

USER intelagent
CMD ["/cmd/intelagent"]

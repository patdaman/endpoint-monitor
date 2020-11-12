# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.15.4

# Copy the local package files to the container's workspace.
# ADD . C://go/src/github.com/patdaman/endpoint-monitor
WORKDIR /go/src/endpoint-monitor
COPY . /go/src/endpoint-monitor

# Get / Update command inside the container.
# RUN go get github.com/patdaman/endpoint-monitor@development

# Build the outyet command inside the container.
RUN go build -o /go/bin/
RUN go install

# RUN cd /go/bin
RUN chmod +x /go/bin/endpoint-monitor

ARG CONFIG_FILE=quest-monitoring.json
ENTRYPOINT [ "endpoint-monitor"]
# CMD ["-config", CONFIG_FILE]
CMD ["-config", "quest-monitoring.json"]

# Document that the service listens 
EXPOSE 80 8083 7321

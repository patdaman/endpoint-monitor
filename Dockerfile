# !!! Do not use	!!!
#This is not ready

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
# ADD . /go/src/github.com/patdaman/endpoint-monitor
ADD . /

# Build the outyet command inside the container.
RUN go get github.com/patdaman/endpoint-monitor@development
# RUN go install https://github.com/patdaman/endpoint-monitor@development
# RUN go install

RUN go build -o main /github.com/patdaman/endpoint-monitor@development

# ENTRYPOINT /go/bin/endpoint-monitor --config /go/src/github.com/patdaman/endpoint-monitor/config.json
ENTRYPOINT /go/bin/endpoint-monitor --config ./quest_monitoring.json

# Document that the service listens 
EXPOSE 80 8083 8086 7321 3000

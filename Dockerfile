FROM golang:1.14 as build

WORKDIR /go/src/github.com/webdevops/pagerduty2elasticsearch-exporter

# Get deps (cached)
COPY ./go.mod /go/src/github.com/webdevops/pagerduty2elasticsearch-exporter
COPY ./go.sum /go/src/github.com/webdevops/pagerduty2elasticsearch-exporter
RUN go mod download

# Compile
COPY ./ /go/src/github.com/webdevops/pagerduty2elasticsearch-exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /pagerduty2elasticsearch-exporter \
    && chmod +x /pagerduty2elasticsearch-exporter
RUN /pagerduty2elasticsearch-exporter --help

#############################################
# FINAL IMAGE
#############################################
FROM gcr.io/distroless/static
COPY --from=build /pagerduty2elasticsearch-exporter /
USER 1000
ENTRYPOINT ["/pagerduty2elasticsearch-exporter"]

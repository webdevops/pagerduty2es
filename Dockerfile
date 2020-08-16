FROM golang:1.15 as build

WORKDIR /go/src/github.com/webdevops/pagerduty2es

# Get deps (cached)
COPY ./go.mod /go/src/github.com/webdevops/pagerduty2es
COPY ./go.sum /go/src/github.com/webdevops/pagerduty2es
COPY ./Makefile /go/src/github.com/webdevops/pagerduty2es
RUN make dependencies

# Compile
COPY ./ /go/src/github.com/webdevops/pagerduty2es
RUN make lint
RUN make build
RUN ./pagerduty2es --help

#############################################
# FINAL IMAGE
#############################################
FROM gcr.io/distroless/static
COPY --from=build /go/src/github.com/webdevops/pagerduty2es/pagerduty2es /
USER 1000
ENTRYPOINT ["/pagerduty2es"]

FROM golang:1.17 as build-env
WORKDIR /go/src/app
ADD main.go /go/src/app
ADD go.mod /go/src/app
ADD go.sum /go/src/app
RUN go mod tidy
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base
LABEL org.opencontainers.image.source="https://github.com/enriquetejeda/github-webhook-docker-stop"
LABEL org.opencontainers.image.authors="Enrique Tejeda"
LABEL org.opencontainers.image.created="2019-10-01"
LABEL org.opencontainers.image.title="github-webhook-docker-stop"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.description="A solution for interact with docker api local with a simple server for listen events from github webhook."
LABEL org.opencontainers.image.licenses="Apache-2.0"
COPY --from=build-env /go/bin/app /
CMD ["/app"]
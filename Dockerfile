FROM golang:1.19

RUN go install github.com/jstemmer/go-junit-report@latest

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/cyverse-de/requests
COPY . .
RUN make

FROM scratch

WORKDIR /app

COPY --from=0 /go/src/github.com/cyverse-de/requests/requests /bin/requests
COPY --from=0 /go/src/github.com/cyverse-de/requests/swagger.json swagger.json

ENTRYPOINT ["requests"]

EXPOSE 8080

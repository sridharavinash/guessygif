FROM golang:1.13.8 as guessy-server
WORKDIR /go/src/github.com/sridharavinash/guessygif
COPY . .
RUN go build -o bin/server

FROM debian:stretch
EXPOSE 8081
COPY --from=guessy-server /go/src/github.com/sridharavinash/guessygif/bin/server /
COPY --from=guessy-server /go/src/github.com/sridharavinash/guessygif/assets /assets
COPY --from=guessy-server /go/src/github.com/sridharavinash/guessygif/public /public

ENTRYPOINT ["/server"]

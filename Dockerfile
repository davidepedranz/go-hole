FROM golang:1.13.1-stretch as builder
WORKDIR $GOPATH/src/github.com/davidepedranz/go-hole/
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-hole

FROM scratch
COPY --from=builder /go-hole /go-hole
COPY /data/blacklist.txt /data/blacklist.txt
ENTRYPOINT ["/go-hole"]
EXPOSE 53/udp

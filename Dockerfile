FROM golang:1.11.1-stretch as builder
RUN wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/bin/dep
WORKDIR $GOPATH/src/github.com/davidepedranz/go-hole/
COPY . ./
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-hole

FROM scratch
COPY --from=builder /go-hole /go-hole
COPY /data/blacklist.txt /data/blacklist.txt
ENTRYPOINT ["/go-hole"]
EXPOSE 53/udp

# go-hole
[![Build Status](https://travis-ci.org/davidepedranz/go-hole.svg?branch=master)](https://travis-ci.org/davidepedranz/go-hole)

`go-hole` is a fast and lightweight [DNS sinkhole](https://en.wikipedia.org/wiki/DNS_Sinkhole) that blocks domains known to serve ads, tracking scripts, malware and other unwanted content. It also caches DNS responses to reduce latency, and collects anonymous statistics about the DNS traffic. `go-hole` is written in Go and runs on every platform and operating systems supported by the Go compiler. `go-hole` can be combined with a private VPN to protect mobile devices on every network.

## TL;TR

Run as a Docker container and use as your primary DNS server:
```sh
docker run --name go-hole -d -p 127.0.0.1:53:8053/udp davidepedranz/go-hole:latest
```

Test that `go-hole` is working correctly:
```sh
nslookup -port=8053 example.com localhost
nslookup -port=8053 googleadservices.com localhost
```

## How does it work?

`go-hole` runs a custom DNS server that selectively blocks unwanted domains by replying `NXDomain (Non-Existent Domain)` to the client. It uses an upstream DNS (by default [1.1.1.1](https://1.1.1.1/)) to resolve the queries the first time, then it caches the response to speed up the following requests.

## Why?

The amount of intrusive ads and tracking services on the Internet is huge and continues to grow. While it is quite easy to block them on a computer using your favourite ad-block plugin, it is difficult or even impossible to do the same on mobile devices. This project aims to block unwanted ads and services at the network level, without the need to install any software on the user's device.

This project is inspired by [Pi-Hole](https://github.com/pi-hole/pi-hole), but with a slightly different approach. `go-hole` provides a single binary that only selectively filters the unwanted domains. The blacklist is static and is loaded at startup and cached in memory.

## Build & Run

```sh
# install the dependencies
dep ensure

# build the binary
go build

# execute the binary
# please make sure the blacklist is available at ./data/blacklist.txt
./go-hole
```

## Configuration

`go-hole` can be configured using a few environment variables:

| Environment Variable | Default Value | Function                                                     |
| -------------------- | ------------- | ------------------------------------------------------------ |
| `DNS_PORT`           | `8053`        | UDP port where to listen for DNS queries.                    |
| `PROMETHEUS_PORT`    | `9090`        | TCP port where to serve the collected metrics.               |
| `UPSTREAM_DNS`       | `1.1.1.1:53`  | IP and port of the upstream DNS to use to resolve the queries. |
| `DEBUG`              | `false`       | If true, `go-hole` logs all queries to the standard output.  |

You can customize the behaviour of `go-hole` by changing domains in the [blacklist](./data/blacklist.txt). The default blacklist can be build with:

```sh
./scripts/make-blacklist.sh
```

## FAQ

### Do you have a Docker container?

Sure, checkout the automatic build on Docker Hub: https://hub.docker.com/r/davidepedranz/go-hole/

### Can I combine it with a VPN software?

Sure, this is the main setup of `go-hole`. For example, you can combine it with [OpenVPN](https://openvpn.net/). We will publish soon a guide to setup `go-hole` and OpenVPN together on a private server.

### Privacy Issues

By default, `go-hole` does not log any DNS query. Logging can be enabled for debug purposes, but we discourage it in production, since it breaches the privacy of the users. On the other hand, `go-hole` is fully instrumented to collect anonymous data about the amount of blocked queries, the response times and other performance metrics.

### Metrics

`go-hole` is instrumented with [Prometheus](https://prometheus.io/) to collect the following metrics:

| Type      | Name                                       | Help                                          |
| --------- | ------------------------------------------ | --------------------------------------------- |
| Histogram | `gohole_dns_queries_duration_seconds`      | Duration of replies to DNS queries.           |
| Histogram | `gohole_blacklist_lookup_duration_seconds` | Duration of a domain lookup in the blacklist. |
| Histogram | `gohole_cache_operation_duration_seconds`  | Duration of an operation on the cache.        |

By default, metrics are served over HTTP at port `9090` and path `/metrics`.

## License

`go-hole` is free software released under the MIT Licence. Please checkout the [LICENSE](./LICENSE) file for details.